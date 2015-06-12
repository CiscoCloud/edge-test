package main

import (
	"flag"

	"github.com/gorilla/websocket"
)

var appPort = flag.Int("port", 9090, "Port to serve on")
var cassandraHost = flag.String("cassandra.host", "localhost", "Cassandra host")
var topics = flag.String("topics", "", "Kafka topics to read from (coma-separated)")
var zkConnect = flag.String("zookeeper", "localhost:2181", "Zookeeper host:port")
var schemaRegistryUrl = flag.String("schema.registry.url", "http://localhost:8081", "Schema registry URL")

type App struct {
	eventFetcher *EventFetcher
	connections  map[*websocket.Conn]chan *Event
}

func NewApp() *App {
	app := new(App)
	config := &EventFetcherConfig{*cassandraHost, *zkConnect, *schemaRegistryUrl, *topics}
	app.eventFetcher = NewEventFetcher(config)
	app.connections = make(map[*websocket.Conn]chan *Event)
	return app
}

func main() {
	flag.Parse()
	app := NewApp()
	defer app.eventFetcher.Close()
	go app.eventFetcher.startFetch()
	app.setHandlers()
	go app.eventSender()
	startWebServer(*appPort)
}
