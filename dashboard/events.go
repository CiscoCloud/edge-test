package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	avro "github.com/stealthly/go-avro"
	kafka "github.com/stealthly/go_kafka_client"
)

type Event struct {
	EventName string `json:"eventName"`
	Second    int64  `json:"second"`
	Operation string `json:"operation"`
	Value     int64  `json:"value"`
	Cnt       int64  `json:"count"`
}

type EventFetcher struct {
	events     chan *Event
	connection *gocql.Session
	config     *EventFetcherConfig
	consumer   *kafka.Consumer
}

func NewEventFetcher(config *EventFetcherConfig) *EventFetcher {
	var err error
	fetcher := new(EventFetcher)
	fetcher.config = config
	fetcher.events = make(chan *Event)
	cluster := gocql.NewCluster(config.CassandraHost)
	cluster.Keyspace = "spark_analysis"
	fetcher.connection, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	fetcher.consumer = fetcher.createConsumer()

	return fetcher
}

func (this *EventFetcher) Close() {
	this.connection.Close()
}

func (this *EventFetcher) createConsumer() *kafka.Consumer {
	coordinatorConfig := kafka.NewZookeeperConfig()
	coordinatorConfig.ZookeeperConnect = []string{this.config.ZkConnect}
	coordinator := kafka.NewZookeeperCoordinator(coordinatorConfig)
	consumerConfig := kafka.DefaultConsumerConfig()
	consumerConfig.AutoOffsetReset = kafka.LargestOffset
	consumerConfig.Coordinator = coordinator
	consumerConfig.Groupid = "event-dashboard"
	consumerConfig.ValueDecoder = kafka.NewKafkaAvroDecoder(this.config.SchemaRegistryUrl)
	consumerConfig.WorkerFailureCallback = func(_ *kafka.WorkerManager) kafka.FailedDecision {
		return kafka.CommitOffsetAndContinue
	}
	consumerConfig.WorkerFailedAttemptCallback = func(_ *kafka.Task, _ kafka.WorkerResult) kafka.FailedDecision {
		return kafka.CommitOffsetAndContinue
	}
	consumerConfig.Strategy = func(_ *kafka.Worker, msg *kafka.Message, taskId kafka.TaskId) kafka.WorkerResult {
		if record, ok := msg.DecodedValue.(*avro.GenericRecord); ok {
			this.events <- &Event{
				EventName: record.Get("eventname").(string),
				Second:    record.Get("second").(int64),
				Operation: record.Get("operation").(string),
				Value:     record.Get("value").(int64),
				Cnt:       record.Get("cnt").(int64),
			}
		} else {
			return kafka.NewProcessingFailedResult(taskId)
		}

		return kafka.NewSuccessfulResult(taskId)
	}

	return kafka.NewConsumer(consumerConfig)
}

func (this *EventFetcher) EventHistory() []Event {
	operations := "'avg10second', 'avg30second', 'avg1minute', 'avg5minute', 'avg10minute', 'avg15minute'"
	timestamp := time.Now().Add(-1 * time.Hour).Unix()
	query := fmt.Sprintf("SELECT second, eventname, operation, cnt, value FROM events WHERE operation IN (%s) AND second > %d;", operations, timestamp)
	data := this.connection.Query(query).Iter()
	var events []Event
	var eventName string
	var second int64
	var operation string
	var value int64
	var cnt int64
	for data.Scan(&second, &eventName, &operation, &cnt, &value) {
		event := Event{
			Second:    second,
			EventName: eventName,
			Operation: operation,
			Cnt:       cnt,
			Value:     value,
		}
		events = append(events, event)
	}
	if err := data.Close(); err != nil {
		log.Fatal(err)
	}
	return events
}

func (this *EventFetcher) startFetch() {
	topicCount := make(map[string]int)
	topicCount[fmt.Sprintf("%s-latencies", this.config.Topic)] = 1
	this.consumer.StartStatic(topicCount)
}

type EventFetcherConfig struct {
	CassandraHost     string
	ZkConnect         string
	SchemaRegistryUrl string
	Topic             string
}
