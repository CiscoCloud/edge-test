edge-test Spark latencies calculation job
================

Pre-Requisites
==============

- java 1.6 or higher   
- Apache Mesos 0.19(or newer) if running on Mesos or Spark master
- Local Spark installation
 
Build Instructions
=================

- Get the project   
```
$ git clone https://github.com/CiscoCloud/edge-test.git
$ cd edge-test
```  
- Build spark job artifact  
```
$ ./gradlew latencies-calc:clean latencies-calc:jar
```

Running on Mesos
================

```
$ ${SPARK_HOME}/bin/spark-submit --master mesos://${MASTER_URL} --spark.mesos.coarse (true|false) --spark.executor.uri ${SPARK_EXECUTOR_URI} --class ly.stealth.latencies.Main latencies-calc/build/libs/latencies-calc-1.0.jar --topics ${TOPICS} --zookeeper ${ZOOKEEPER} --broker.list ${BROKERS} --schema.registry.url ${SCHEMA_REGISTRY_URL} --partitions ${SPARK_PARTITIONS}
```

Running on Standalone Spark Master
==================================

```
$ ${SPARK_HOME}/bin/spark-submit --master spark://${MASTER_URL} --class ly.stealth.latencies.Main latencies-calc/build/libs/latencies-calc-1.0.jar --topics ${TOPICS} --zookeeper ${ZOOKEEPER} --broker.list ${BROKERS} --schema.registry.url ${SCHEMA_REGISTRY_URL} --partitions ${SPARK_PARTITIONS}
```

*List of available flags:*

```
--topics: Comma separated list of topics to read data from.
--zookeeper: Zookeeper connection string - host:port.
--broker.list: Comma separated string of host:port.
--schema.registry.url: Schema registry URL starting with http:// or https://.
--partitions: Initial amount of RDD partitions.
```