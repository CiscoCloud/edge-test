package main

import (
	"fmt"
	"github.com/gocql/gocql"
	avro "github.com/stealthly/go-avro"
	kafka "github.com/stealthly/go_kafka_client"
	"log"
	"reflect"
	"strings"
	"sync"
)

var REQUIRED_KEYS = []string{"datacenter", "floor", "tile", "rack"}

type AvroCassandraConsumerConfig struct {
	CassandraHost     string
	CassandraKeyspace string
	ZkConnect         []string
	SchemaRegistryUrl string
	Topics            []string
	NumStreams        int
	Group             string
}

type AvroCassandraConsumer struct {
	config           *AvroCassandraConsumerConfig
	consumer         *kafka.Consumer
	cassandraSession *gocql.Session
	batch            *gocql.Batch
	batchLock        sync.RWMutex
}

func NewAvroCassandraConsumer(config *AvroCassandraConsumerConfig) *AvroCassandraConsumer {
	acConsumer := new(AvroCassandraConsumer)
	acConsumer.config = config

	coordinatorConfig := kafka.NewZookeeperConfig()
	coordinatorConfig.ZookeeperConnect = config.ZkConnect
	coordinator := kafka.NewZookeeperCoordinator(coordinatorConfig)

	consumerConfig := kafka.DefaultConsumerConfig()
	consumerConfig.Groupid = config.Group
	consumerConfig.Coordinator = coordinator
	consumerConfig.ValueDecoder = kafka.NewKafkaAvroDecoder(config.SchemaRegistryUrl)
	consumerConfig.Strategy = func(worker *kafka.Worker, message *kafka.Message, taskId kafka.TaskId) kafka.WorkerResult {
		if decodedMessage, ok := message.DecodedValue.(*avro.GenericRecord); ok {
			keyTags := getKeyTags(decodedMessage.Get("tag").(map[string]interface{}))
			if len(keyTags) != len(REQUIRED_KEYS) {
				kafka.Errorf(acConsumer, "Invalid message: required fields missing(have: %v, wanted: %v)", keyTags, REQUIRED_KEYS)
				return kafka.NewProcessingFailedResult(taskId)
			}

			compositeKey := make(map[string]string)
			keys := make([]string, 0)
			values := make([]string, 0)
			for _, key := range REQUIRED_KEYS {
				keys = append(keys, key)
				values = append(values, keyTags[key].(string))
				compositeKey[strings.Join(keys, "_")] = strings.Join(values, "|")
			}

			fields := decodedMessage.Schema().(*avro.RecordSchema).Fields
			updateValues := extractValues(decodedMessage)
			updateClause := strings.Join(updateValues, ",")
			for tableSuffix, id := range compositeKey {
				insertQuery := fmt.Sprintf("UPDATE events_by_%s SET %v WHERE id = '%s' and time = dateof(now())", tableSuffix, updateClause, id)
				err := acConsumer.cassandraSession.Query(insertQuery).Exec()
				if err != nil {
					kafka.Warnf(acConsumer, "Table events_by_%s does not exist yet. Trying to create one...", tableSuffix)
					fieldMappings := make([]string, 0)
					for _, field := range fields {
						fieldMappings = append(fieldMappings, fmt.Sprintf("%s %s", field.Name, mapToCassandraType(field.Type)))
					}

					createQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS events_by_%s(id text, time bigint,  %s, PRIMARY KEY(id, time)) WITH CLUSTERING ORDER BY (time DESC)",
						tableSuffix, strings.Join(fieldMappings, ","))
					if err = acConsumer.cassandraSession.Query(createQuery).Exec(); err != nil {
						kafka.Errorf(acConsumer, "Error executing query %s due to: %s", createQuery, err.Error())
						return kafka.NewProcessingFailedResult(taskId)
					}
					kafka.Infof(acConsumer, "Successfully created events_by_%s table", tableSuffix)

					if err = acConsumer.cassandraSession.Query(insertQuery).Exec(); err != nil {
						kafka.Errorf(acConsumer, "Error executing query %s due to: %s", insertQuery, err.Error())
						return kafka.NewProcessingFailedResult(taskId)
					}
				}
			}

			return kafka.NewSuccessfulResult(taskId)
		}

		return kafka.NewProcessingFailedResult(taskId)
	}
	consumerConfig.WorkerFailureCallback = func(wm *kafka.WorkerManager) kafka.FailedDecision {
		kafka.Error(acConsumer, "Failed to write critical amount of messages into Cassandra. Shutting down...")
		return kafka.DoNotCommitOffsetAndStop
	}
	consumerConfig.WorkerFailedAttemptCallback = func(task *kafka.Task, result kafka.WorkerResult) kafka.FailedDecision {
		kafka.Errorf(acConsumer, "Failed to write %v to the Cassandra after %d retries", task.Msg.DecodedValue, task.Retries)
		return kafka.DoNotCommitOffsetAndContinue
	}
	acConsumer.consumer = kafka.NewConsumer(consumerConfig)

	cluster := gocql.NewCluster(config.CassandraHost)
	cluster.Keyspace = config.CassandraKeyspace
	var err error
	acConsumer.cassandraSession, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %s", err.Error())
	}
	acConsumer.batch = acConsumer.cassandraSession.NewBatch(gocql.LoggedBatch)

	return acConsumer
}

func (this *AvroCassandraConsumer) Start() {
	topicMap := make(map[string]int)
	for _, topic := range this.config.Topics {
		topicMap[topic] = this.config.NumStreams
	}
	this.consumer.StartStatic(topicMap)
	this.cassandraSession.Close()
}

func (this *AvroCassandraConsumer) Stop() <-chan bool {
	return this.consumer.Close()
}

func (this *AvroCassandraConsumer) String() string {
	return fmt.Sprintf("ac-consumer")
}

func getKeyTags(tags map[string]interface{}) map[string]interface{} {
	keyTags := make(map[string]interface{})
	for _, key := range REQUIRED_KEYS {
		if value, ok := tags[key]; ok {
			keyTags[key] = value
		}
	}

	return keyTags
}

func mapToCassandraType(fieldType avro.Schema) string {
	switch fieldType.Type() {
	case avro.Array:
		return fmt.Sprintf("list<%s>", mapToCassandraType(fieldType.(*avro.ArraySchema).Items))
	case avro.Map:
		return fmt.Sprintf("map<text, %s>", mapToCassandraType(fieldType.(*avro.MapSchema).Values))
	case avro.String:
		return "text"
	case avro.Bytes:
		return "blob"
	case avro.Int:
		return "int"
	case avro.Long:
		return "bigint"
	case avro.Float:
		return "float"
	case avro.Double:
		return "double"
	case avro.Boolean:
		return "boolean"
	case avro.Union:
		return mapToCassandraType(fieldType.(*avro.UnionSchema).Types[1])
	case avro.Record:
		result := make([]string, 0)
		for _, field := range fieldType.(*avro.RecordSchema).Fields {
			result = append(result, mapToCassandraType(field.Type))
		}
		return fmt.Sprintf("frozen<tuple<%s>>", strings.Join(result, ", "))
	}

	panic(fmt.Sprintf("Unsupported type: %s", fieldType.GetName()))
}

func extractValues(record *avro.GenericRecord) []string {
	extractedValues := make([]string, 0)
	fields := record.Schema().(*avro.RecordSchema).Fields
	for _, field := range fields {
		extractedValues = append(extractedValues, fmt.Sprintf("%s = %s", field.Name, mapToCassandraValue(record.Get(field.Name))))
	}

	return extractedValues
}

func mapToCassandraValue(obj interface{}) string {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	switch v.Kind() {
	case reflect.Ptr: {
		if record, ok := v.Elem().Interface().(avro.GenericRecord); ok {
			result := make([]string, 0)
			fields := record.Schema().(*avro.RecordSchema).Fields
			for _, field := range fields {
				result = append(result, mapToCassandraValue(record.Get(field.Name)))
			}

			return fmt.Sprintf("(%s)", strings.Join(result, ", "))
		} else {
			return mapToCassandraValue(v.Elem().Interface())
		}
	}
	case reflect.String:
		return fmt.Sprintf("'%v'", v.Interface())
	case reflect.Map:
		{
			result := make([]string, v.Len())
			keys := v.MapKeys()
			for i := 0; i < v.Len(); i++ {
				result[i] = fmt.Sprintf("%s: %s", mapToCassandraValue(keys[i].Interface()),
					mapToCassandraValue(v.MapIndex(keys[i]).Interface()))
			}

			return fmt.Sprintf("{%s}", strings.Join(result, ", "))
		}
	case reflect.Array | reflect.Slice:
		{
			result := make([]string, v.Len())
			for i := 0; i < v.Len(); i++ {
				result[i] = mapToCassandraValue(v.Index(i).Interface())
			}

			return fmt.Sprintf("[%s]", strings.Join(result, ", "))
		}
	case reflect.Struct:
		{
			result := make([]string, t.NumField())
			for i := 0; i < t.NumField(); i++ {
				result[i] = mapToCassandraValue(v.Field(i).Interface())
			}

			return fmt.Sprintf("(%s)", strings.Join(result, ", "))
		}
	default:{
		return fmt.Sprintf("%v", v.Interface())
	}
	}
}
