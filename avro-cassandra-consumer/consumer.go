package main

import (
	"fmt"
	"github.com/gocql/gocql"
	kafka "github.com/stealthly/go_kafka_client"
	"log"
	"sync"
	"time"
)

type AvroCassandraConsumerConfig struct {
	CassandraHost     string
	CassandraKeyspace string
	ZkConnect         []string
	SchemaRegistryUrl string
	Topics            []string
	NumStreams        int
	Group             string
	UpdateInterval    time.Duration
}

type AvroCassandraConsumer struct {
	config           *AvroCassandraConsumerConfig
	consumer         *kafka.Consumer
	cassandraSession *gocql.Session
	batch            *gocql.Batch
	batchLock        sync.RWMutex
	shuttingDown     bool
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
		inReadLock(&acConsumer.batchLock, func() {
			//TODO: finish up the query
			acConsumer.batch.Query("UPDATE events where ...")
		})
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

	go func() {
		for !this.shuttingDown {
			<-time.After(this.config.UpdateInterval)
			inWriteLock(&this.batchLock, func() {
				if this.batch.Size() == 0 {
					kafka.Debug(this, "Batch is empty")
					return
				}

				kafka.Debug(this, "Trying to update Cassandra")
				if err := this.cassandraSession.ExecuteBatch(this.batch); err != nil {
					log.Fatalf("Failed to update Cassandra: %s", err.Error())
				}
				kafka.Debugf(this, "Successfully updated Cassandra with %d values", this.batch.Size())

				this.batch = this.cassandraSession.NewBatch(gocql.LoggedBatch)
			})
		}
	}()
	this.consumer.StartStatic(topicMap)
}

func (this *AvroCassandraConsumer) Stop() <-chan bool {
	this.shuttingDown = true
	this.cassandraSession.Close()
	return this.consumer.Close()
}

func (this *AvroCassandraConsumer) String() string {
	return fmt.Sprintf("ac-consumer")
}

func inReadLock(lock *sync.RWMutex, fun func()) {
	lock.RLock()
	defer lock.RUnlock()

	fun()
}

func inWriteLock(lock *sync.RWMutex, fun func()) {
	lock.Lock()
	defer lock.Unlock()

	fun()
}
