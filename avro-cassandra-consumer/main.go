package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var cassandraHost = flag.String("cassandra.host", "localhost", "Cassandra host")
var cassandraKeyspace = flag.String("cassandra.keyspace", "spark_analysis", "Cassandra keyspace")
var topics = flag.String("topics", "", "Comma separated kafka topics to read from")
var numStreams = flag.Int("num.streams", 1, "Amount of streams for each topic")
var group = flag.String("group", "ac-consumer", "Kafka consumer group name")
var zkConnect = flag.String("zookeeper", "localhost:2181", "Zookeeper host:port")
var schemaRegistryUrl = flag.String("schema.registry.url", "http://localhost:8081", "Schema registry URL")
var updateInterval = flag.String("update.interval", "1s", "Interval at which Cassandra should be updated")

func main() {
	flag.Parse()
	if *topics == "" {
		log.Fatal("You have to provide at least one topic")
	}

	parsedInterval, err := time.ParseDuration(*updateInterval)
	if err != nil {
		log.Fatalf("Incorrect update interval argument: %s", *updateInterval)
	}

	config := AvroCassandraConsumerConfig{
		CassandraHost:     *cassandraHost,
		CassandraKeyspace: *cassandraKeyspace,
		Topics:            strings.Split(*topics, ","),
		NumStreams:        *numStreams,
		Group:             *group,
		ZkConnect:         strings.Split(*zkConnect, ","),
		SchemaRegistryUrl: *schemaRegistryUrl,
		UpdateInterval:    parsedInterval,
	}
	acConsumer := NewAvroCassandraConsumer(&config)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	go func() {
		<-ctrlc
		<-acConsumer.Stop()
	}()

	acConsumer.Start()
}
