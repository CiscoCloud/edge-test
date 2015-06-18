package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
)

var cassandraHost = flag.String("cassandra.host", "localhost", "Cassandra host")
var cassandraKeyspace = flag.String("cassandra.keyspace", "spark_analysis", "Cassandra keyspace")
var topics = flag.String("topic", "", "Comma separated kafka topics to read from")
var numStreams = flag.Int("num.streams", 1, "Amount of streams for each topic")
var group = flag.String("group", "ac-consumer", "Kafka consumer group name")
var zkConnect = flag.String("zookeeper", "localhost:2181", "Zookeeper host:port")
var schemaRegistryUrl = flag.String("schema.registry.url", "http://localhost:8081", "Schema registry URL")

func main() {
	flag.Parse()
	config := AvroCassandraConsumerConfig{
		CassandraHost:     *cassandraHost,
		CassandraKeyspace: *cassandraKeyspace,
		Topics:            strings.Split(*topics, ","),
		NumStreams:        *numStreams,
		Group:             *group,
		ZkConnect:         strings.Split(*zkConnect, ","),
		SchemaRegistryUrl: *schemaRegistryUrl,
	}
	acConsumer := NewAvroCassandraConsumer(config)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	go func() {
		<-ctrlc
		<-acConsumer.Stop()
	}()

	acConsumer.Start()
}
