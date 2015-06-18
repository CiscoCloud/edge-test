package main

import (
	"fmt"
	"github.com/gocql/gocql"
	kafka "github.com/stealthly/go_kafka_client"
	"os"
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
		//Writing to Cassandra should go here
		return kafka.NewSuccessfulResult(taskId)
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
		kafka.Error(acConsumer, err)
		os.Exit(1)
	}

	return acConsumer
}

func (this *AvroCassandraConsumer) Start() {
	topicMap := make(map[string]int)
	for _, topic := range this.config.Topics {
		topicMap[topic] = this.config.NumStreams
	}
	this.consumer.StartStatic(topicMap)
}

func (this *AvroCassandraConsumer) Stop() <-chan bool {
	this.cassandraSession.Close()
	return this.Stop()
}

func (this *AvroCassandraConsumer) String() string {
	return fmt.Sprintf("ac-consumer")
}
