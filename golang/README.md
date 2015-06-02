edge-test golang
================

Pre-Requisites
==============

- [Golang](http://golang.org/doc/install)   
- A standard and working Go workspace setup   
- [godep](https://github.com/tools/godep)   
- Apache Mesos 0.19 or newer
 
Build Instructions
=================

- Get the project   
```
$ cd $GOPATH/src/
$ mkdir -p github.com/CiscoCloud
$ cd github.com/CiscoCloud
$ git clone https://github.com/CiscoCloud/edge-test.git
$ cd edge-test/golang
$ godep restore
```

- Build the scheduler and the executor
```
$ cd $GOPATH/src/github.com/CiscoCloud/edge-test/golang
$ go build -tags framework framework.go
$ go build -tags executor executor.go
```
- Package the executor (**make sure the built binary has executable permissions before this step!**)
```
$ zip -r executor.zip executor
```
- Place the built framework and executor archive somewhere on Mesos Master node

Running
=======

You will need a running Mesos master and slaves to run. The following commands should be launched on Mesos Master node.

```
$ cd <framework-location>
$ ./framework --master master:5050 --producer.config producer.config --topic logs
```

*List of available flags:*

```
--artifact.host="master": Binding host for artifact server.
--artifact.port=6666: Binding port for artifact server.
--cpu.per.task=0.2: CPUs per task.
--executor.archive="executor.zip": Executor archive name. Absolute or relative path are both ok.
--executor.name="executor": Executor binary name contained in archive.
--instances=1: Number of tasks to run.
--master="127.0.0.1:5050": Mesos Master address <ip:port>.
--mem.per.task=256: Memory per task.
--producer.config: Producer properties file name.
--sync: Flag to respond only after decoding-encoding is done.
--topic: Topic to produce transformed data to.
```

Usage
=====

To distinguish various types of data `Content-Type` header is used. The following 3 are supported:

- `application/json` for json data
- `application/x-protobuf` for protobuf data
- `avro/binary` for avro data

Examples of data expected are located in `testdata` directory. You may try them out as follows:

JSON:
```
curl --header "Content-Type: application/json" --data-binary "@$GOPATH/src/github.com/CiscoCloud/edge-test/golang/testdata/test-json" http://master:31000
```

Avro:
```
curl --header "Content-Type: avro/binary" --data-binary "@$GOPATH/src/github.com/CiscoCloud/edge-test/golang/testdata/test-avro" http://master:31000
```

Protobuf:
```
curl --header "Content-Type: application/x-protobuf" --data-binary "@$GOPATH/src/github.com/CiscoCloud/edge-test/golang/testdata/test-proto" http://master:31000
```