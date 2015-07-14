Frameworks Dashboard
====================

Pre-Requisites
==============

- Go 1.4 or higher

Build Instructions
==================

```
$ git clone https://github.com/CiscoCloud/edge-test.git
$ cd edge-test/dashboard
$ go build .
```

Running
=======

```
./dashboard --zookeeper <zkhost>:<zkport> --topic <topicname> --schema.registry.url <schemaregistry_url> --cassandra.host <cassandra_host>
```

*List of available flags:*

```
--port: Port for web interface. Default to 9090.
--cassandra.host: Hostname of ip of cassandra to get metrics from. Default to localhost.
--topic: Kafka topic to read from.
--zookeeper: Zookeeper host:port. Default to localhost:2181.
--schema.registry.url: Schema registry URL starting with http:// or https://. Default to http://localhost:8081
```
