package main

import (
	"fmt"
	"github.com/gocql/gocql"
	avro "github.com/stealthly/go-avro"
	kafka "github.com/stealthly/go_kafka_client"
	"log"
	"strings"
	"sync"
)

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
	consumerConfig.AutoOffsetReset = kafka.SmallestOffset
	consumerConfig.Coordinator = coordinator
	consumerConfig.ValueDecoder = kafka.NewKafkaAvroDecoder(config.SchemaRegistryUrl)
	consumerConfig.Strategy = func(worker *kafka.Worker, message *kafka.Message, taskId kafka.TaskId) kafka.WorkerResult {
		if decodedMessage, ok := message.DecodedValue.(*avro.GenericRecord); ok {
			tags := decodedMessage.Get("tag").(map[string]interface{})
			compositeKey := make(map[string]string)
			keys := make([]string, 0)
			values := make([]string, 0)
			for key, value := range tags {
				keys = append(keys, key)
				values = append(values, value.(string))
				compositeKey[strings.Join(keys, "_")] = strings.Join(values, "|")
			}

			updateFieldBits := make([]string, 0)
			updateValues := make([]interface{}, 0)
			fields := decodedMessage.Schema().(*avro.RecordSchema).Fields
			for _, field := range fields {
				updateFieldBits = append(updateFieldBits, fmt.Sprintf("%s = ?", field.Name))
				updateValues = append(updateValues, decodedMessage.Get(field.Name))
			}
			updateClause := strings.Join(updateFieldBits, ",")
			for tableSuffix, id := range compositeKey {
				insertQuery := fmt.Sprintf("UPDATE events_by_%s SET %s WHERE id = '%s' and time = dateof(now())", tableSuffix, updateClause, id)
				err := acConsumer.cassandraSession.Query(insertQuery, updateValues...).Exec()
				if err != nil {
					kafka.Warnf(acConsumer, "Table events_by_%s does not exist yet. Trying to create one...", tableSuffix)
					fieldMappings := make([]string, 0)
					for _, field := range fields {
						fieldMappings = append(fieldMappings, fmt.Sprintf("%s %s", field.Name, mapType(field.Type)))
					}

					createQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS events_by_%s(id text, time bigint,  %s, PRIMARY KEY(id, time)) WITH CLUSTERING ORDER BY (time DESC)",
						tableSuffix, strings.Join(fieldMappings, ","))
					if err = acConsumer.cassandraSession.Query(createQuery).Exec(); err != nil {
						kafka.Errorf(acConsumer, "Error executing query %s due to: %s", createQuery, err.Error())
						return kafka.NewProcessingFailedResult(taskId)
					}
					kafka.Infof(acConsumer, "Successfully created events_by_%s table", tableSuffix)

					if err = acConsumer.cassandraSession.Query(insertQuery, updateValues...).Exec(); err != nil {
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

func mapType(fieldType avro.Schema) string {
	switch fieldType.Type() {
	case avro.Array:
		return fmt.Sprintf("list<%s>", mapType(fieldType.(*avro.ArraySchema).Items))
	case avro.Map:
		return fmt.Sprintf("map<text, %s>", mapType(fieldType.(*avro.MapSchema).Values))
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
		return mapType(fieldType.(*avro.UnionSchema).Types[1])
	}

	panic(fmt.Sprintf("Unknown type: %s", fieldType.GetName()))
}

func (this *AvroCassandraConsumer) Stop() <-chan bool {
	return this.consumer.Close()
}

func (this *AvroCassandraConsumer) String() string {
	return fmt.Sprintf("ac-consumer")
}
