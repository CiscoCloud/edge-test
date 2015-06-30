package main

import (
	"flag"
	kafka "github.com/stealthly/go_kafka_client"
	"log"
	"os"
	"os/signal"
	"strings"
)

var flags = flag.NewFlagSet("Kafka Avro Consumer", flag.ExitOnError)
var cassandraHost = flags.String("cassandra.host", "", "Cassandra host --cassandra.host [host]")
var cassandraKeyspace = flags.String("cassandra.keyspace", "spark_analysis", "Cassandra keyspace")
var topics = flags.String("topics", "", "Comma separated list of kafka topics: --topics [topic1,topic2...]")
var numStreams = flags.Int("num.streams", 1, "Amount of streams for each topic")
var group = flags.String("group", "ac-consumer", "Kafka consumer group name")
var zkConnect = flags.String("zookeeper", "", "Zookeeper connection string: --zookeeper [host:port]")
var schemaRegistryUrl = flags.String("schema.registry.url", "", "Avro Schema Registry URL: --schema.registry.url http://[host:port]")
var logLevel = flags.String("log.level", "info", "Log level")

func main() {
	flags.Parse(os.Args[1:])
	if logLevel != nil {
		setLogLevel(*logLevel)
	}
	if *cassandraHost == "" {
		log.Fatal("You have to provide Cassandra host: --cassandra.host [host]")
	}
	if *zkConnect == "" {
		log.Fatal("You have to provide at least one zookeeper host: --zookeeper [host:port]")
	}
	if *schemaRegistryUrl == "" {
		log.Fatal("You have to provide Schema Registry URL: --schema.registry.url http://[host:port]")
	}
	if *topics == "" {
		log.Fatal("You have to provide at least one topic: --topics [topic1,topic2...]")
	}

	config := AvroCassandraConsumerConfig{
		CassandraHost:     *cassandraHost,
		CassandraKeyspace: *cassandraKeyspace,
		Topics:            strings.Split(*topics, ","),
		NumStreams:        *numStreams,
		Group:             *group,
		ZkConnect:         strings.Split(*zkConnect, ","),
		SchemaRegistryUrl: *schemaRegistryUrl,
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

func setLogLevel(logLevel string) {
	var level kafka.LogLevel
	switch strings.ToLower(logLevel) {
	case "trace":
		level = kafka.TraceLevel
	case "debug":
		level = kafka.DebugLevel
	case "info":
		level = kafka.InfoLevel
	case "warn":
		level = kafka.WarnLevel
	case "error":
		level = kafka.ErrorLevel
	case "critical":
		level = kafka.CriticalLevel
	default:
		log.Fatalf("Invalid log level: %s", logLevel)
	}
	kafka.Logger = kafka.NewDefaultLogger(level)
}
