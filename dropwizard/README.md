edge-test dropwizard
================

Pre-Requisites
==============

- java 1.7 or higher   
- Apache Mesos 0.19 or newer
 
Build Instructions
=================

- Get the project   
```
$ git clone https://github.com/CiscoCloud/edge-test.git
$ cd edge-test/dropwizard
```

- Build the scheduler and the executor
```
$ ./gradlew clean jar
```
- Place the built jar file, `config.yml` and `executor.yml` somewhere on Mesos Master node

Running
=======

You will need a running Mesos master and slaves to run. The following commands should be launched on Mesos Master node.

```
$ cd <framework-location>
$ java -jar logging-mesos-0.1.jar --master master:5050 --user vagrant --executor logging-mesos-0.1.jar --schema.registry http://192.168.3.1:8081 --broker.list 192.168.3.1:9092 --topic logs
```

*List of available flags:*

```
--artifact.host="master": Binding host for artifact server.
--artifact.port=6666: Binding port for artifact server. This should match the `server.applicationConnectors.port` value in `config.yml`
--cpu.per.task=0.2: CPUs per task.
--executor="logging-mesos-0.1.jar": Executor file name. Absolute or relative path are both ok.
--instances=1: Number of tasks to run.
--master="127.0.0.1:5050": Mesos Master address <ip:port>.
--mem.per.task=256: Memory per task.
--schema.registry.url: Avro Schema Registry URL.
--topic: Topic to produce transformed data to.
--user: Mesos user to run for.
```

Usage
=====

To distinguish various types of data `Content-Type` header is used. The following 3 are supported:

- `application/json` for json data
- `application/x-protobuf` for protobuf data
- `avro/binary` for avro data

Usage examples can be found [here](https://github.com/CiscoCloud/edge-test/tree/master/golang#usage)

Scaling
=======

You can scale tasks dynamically by calling `http://<master>:<artifact.port>/scale/<instances>` e.g. `http://master:6666/scale/6` means you want to scale to 6 tasks (no matter if this is upscale or downscale, both should work).

Updating producer properties
============================

You may provide configurations for producers used in executors by adding configurations to `producer.config` file. This changes will take effect when running new tasks, e.g. if you add some properties, all tasks executed after that will see that changes.