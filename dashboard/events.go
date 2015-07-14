package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"
	avro "github.com/stealthly/go-avro"
	kafka "github.com/stealthly/go_kafka_client"
)

type Event struct {
	EventName     string `json:"eventName"`
	Second        int64  `json:"second"`
	Framework     string `json:"framework"`
	Latency       int64  `json:"value"`
	ReceivedCount int64  `json:"count_received"`
	SentCount     int64  `json:"count_sent"`
}

type EventFetcher struct {
	events     chan *Event
	connection *gocql.Session
	config     *EventFetcherConfig
	consumer   *kafka.Consumer
	connected  bool
}

func NewEventFetcher(config *EventFetcherConfig) *EventFetcher {
	var err error
	fetcher := new(EventFetcher)
	fetcher.config = config
	fetcher.events = make(chan *Event)
	cluster := gocql.NewCluster(config.CassandraHost)
	cluster.Keyspace = "spark_analysis"
	fetcher.connection, err = cluster.CreateSession()
	fetcher.connected = true
	if err != nil {
		log.Printf("Can't connect to Cassandra, %q", err)
		fetcher.connected = false
	}
	fetcher.consumer, err = fetcher.createConsumer()
	if err != nil {
		log.Printf("Can't connect to Zookeeper, %q", err)
		fetcher.connected = false
	}

	return fetcher
}

func (this *EventFetcher) Close() {
	this.connection.Close()
}

func (this *EventFetcher) createConsumer() (*kafka.Consumer, error) {
	fmt.Println(this.config.ZkConnect)
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
				EventName:     record.Get("eventname").(string),
				Second:        record.Get("second").(int64),
				Framework:     record.Get("framework").(string),
				Latency:       record.Get("latency").(int64),
				ReceivedCount: record.Get("received_count").(int64),
				SentCount:     record.Get("sent_count").(int64),
			}
		} else {
			return kafka.NewProcessingFailedResult(taskId)
		}

		return kafka.NewSuccessfulResult(taskId)
	}

	consumer, err := kafka.NewSlaveConsumer(consumerConfig)
	return consumer, err
}

func (this *EventFetcher) EventHistory() []Event {
	frameworks := "'Golang', 'Dropwizard', 'Finagle', 'Spray', 'Play', 'Unfiltered', 'Netty'"
	timestamp := time.Now().Add(-1 * time.Hour).Unix()
	query := fmt.Sprintf("SELECT second, eventname, framework, received_count, sent_count, latency FROM events WHERE framework IN (%s) AND second > %d;", frameworks, timestamp)
	data := this.connection.Query(query).Iter()
	var events []Event
	var eventName string
	var second int64
	var framework string
	var latency int64
	var received_count int64
	var sent_count int64
	for data.Scan(&second, &eventName, &framework, &received_count, &sent_count, &latency) {
		event := Event{
			Second:        second,
			EventName:     eventName,
			Framework:     framework,
			ReceivedCount: received_count,
			SentCount:     sent_count,
			Latency:       latency,
		}
		events = append(events, event)
	}
	if err := data.Close(); err != nil {
		log.Fatal(err)
	}
	return events
}

func (this *EventFetcher) startFetch() {
	if !this.connected {
		return
	}
	topicCount := make(map[string]int)
	topics := strings.Split(this.config.Topics, ",")
	for _, topic := range topics {
		topicCount[fmt.Sprintf("%s-latencies", topic)] = 1
	}
	this.consumer.StartStatic(topicCount)
}

type EventFetcherConfig struct {
	CassandraHost     string
	ZkConnect         string
	SchemaRegistryUrl string
	Topics            string
}
