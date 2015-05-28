package ly.stealth.latencies

import java.util.{Properties, UUID}

import _root_.kafka.serializer.DefaultDecoder
import com.datastax.spark.connector._
import com.datastax.spark.connector.cql.CassandraConnector
import io.confluent.kafka.serializers.KafkaAvroSerializer
import org.apache.avro.Schema
import org.apache.avro.generic.GenericData.Record
import org.apache.avro.generic.{GenericData, GenericRecord}
import org.apache.avro.util.Utf8
import org.apache.kafka.clients.producer.ProducerConfig._
import org.apache.kafka.clients.producer.{KafkaProducer, ProducerRecord}
import org.apache.spark.SparkConf
import org.apache.spark.storage.StorageLevel
import org.apache.spark.streaming._
import org.apache.spark.streaming.dstream.DStream
import org.apache.spark.streaming.kafka._

object Main extends App {
  val parser = new scopt.OptionParser[AppConfig]("spark-analysis") {
    head("Latencies calculation job", "1.0")
    opt[String]("topic") unbounded() required() action { (value, config) =>
      config.copy(topic = value)
    } text ("Topic to read data from")
    opt[String]("zookeeper") unbounded() required() action { (value, config) =>
      config.copy(zookeeper = value)
    } text ("Zookeeper connection string - host:port")
    opt[String]("broker.list") unbounded() required() action { (value, config) =>
      config.copy(brokerList = value)
    } text ("Comma separated string of host:port")
    opt[String]("schema.registry.url") unbounded() required() action { (value, config) =>
      config.copy(schemaRegistryUrl = value)
    } text ("Schema registry URL")
    checkConfig { c =>
      if (c.topic.isEmpty || c.brokerList.isEmpty) {
        failure("You haven't provided all required parameters")
      } else {
        success
      }
    }
  }
  val appConfig = parser.parse(args, AppConfig()) match {
    case Some(c) => c
    case None => sys.exit(1)
  }

  val sparkConfig = new SparkConf().setAppName("spark-analysis").set("spark.serializer", "org.apache.spark.serializer.KryoSerializer")
  val ssc = new StreamingContext(sparkConfig, Seconds(1))
  ssc.checkpoint("spark-analysis")

  val cassandraConnector = CassandraConnector(sparkConfig)
  cassandraConnector.withSessionDo(session => {
    session.execute("CREATE KEYSPACE IF NOT EXISTS spark_analysis WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1}")
    session.execute("CREATE TABLE IF NOT EXISTS spark_analysis.events(eventname text, second int, operation text, value int, ntpstatus int, cnt int, PRIMARY KEY(eventname, second, operation))")
  })

  val consumerConfig = Map(
    "group.id" -> "spark-analysis-%s".format(UUID.randomUUID.toString),
    "zookeeper.connect" -> appConfig.zookeeper,
    "auto.offset.reset" -> "largest",
    "schema.registry.url" -> appConfig.schemaRegistryUrl)
  val producerConfig = new Properties()
  producerConfig.put(BOOTSTRAP_SERVERS_CONFIG, appConfig.brokerList)
  producerConfig.put(KEY_SERIALIZER_CLASS_CONFIG, classOf[KafkaAvroSerializer])
  producerConfig.put(VALUE_SERIALIZER_CLASS_CONFIG, classOf[KafkaAvroSerializer])
  producerConfig.put("schema.registry.url", appConfig.schemaRegistryUrl)

  start(ssc, consumerConfig, producerConfig, appConfig.topic)

  ssc.start()
  ssc.awaitTermination()

  def start(ssc: StreamingContext, consumerConfig: Map[String, String], producerConfig: Properties, topic: String) = {
    val stream = KafkaUtils.createStream[Array[Byte], SchemaAndData, DefaultDecoder, AvroDecoder](ssc, consumerConfig, Map(topic -> 2), StorageLevel.MEMORY_AND_DISK_SER).persist()
    stream.persist()
    calculateAverages(stream, "second", 10, topic, producerConfig)
    calculateAverages(stream, "second", 30, topic, producerConfig)
    calculateAverages(stream, "minute", 1, topic, producerConfig)
    calculateAverages(stream, "minute", 5, topic, producerConfig)
    calculateAverages(stream, "minute", 10, topic, producerConfig)
    calculateAverages(stream, "minute", 15, topic, producerConfig)
  }

  def calculateAverages(stream: DStream[(Array[Byte], SchemaAndData)], durationUnit: String, durationValue: Long, topic: String, producerConfig: Properties) = {
    val latencyStream = stream.window(windowDuration(durationUnit, durationValue)).map(value => {
      val record = value._2.deserialize().asInstanceOf[GenericRecord]
      import scala.collection.JavaConversions._
      val timings = record.get("timings").asInstanceOf[GenericData.Array[Record]]
      timings.combinations(2).map(entry => {
        (entry.head.get("key").asInstanceOf[Utf8].toString + "-" + entry.last.get("key").asInstanceOf[Utf8].toString,
          entry.last.get("value").asInstanceOf[Long] - entry.head.get("value").asInstanceOf[Long],
          entry.last.get("ntpstatus").asInstanceOf[Long] - entry.head.get("ntpstatus").asInstanceOf[Long])
      }).toList
    }).reduce((acc, value) => {
      acc ++ value
    }).flatMap(entry => {
      val second = System.currentTimeMillis()/1000
      entry.groupBy(entry => (entry._1)).map { case (key, values) => {
        val timings = values.map(_._2)
        val ntpstatus = values.map(_._3)
        Event(key, second, "avg%d%s".format(durationValue, durationUnit), timings.sum / timings.size, ntpstatus.sum / ntpstatus.size, timings.size)
      }
      }
    }).persist()

    val schema = "{\"type\":\"record\",\"name\":\"event\",\"fields\":[{\"name\":\"eventname\",\"type\":\"string\"},{\"name\":\"second\",\"type\":\"long\"},{\"name\":\"operation\",\"type\":\"string\"},{\"name\":\"value\",\"type\":\"long\"},{\"name\":\"ntpstatus\",\"type\":\"long\"},{\"name\":\"cnt\",\"type\":\"long\"}]}"
    latencyStream.foreachRDD(rdd => {
      rdd.foreachPartition(latencies => {
        val producer = new KafkaProducer[Any, AnyRef](producerConfig)
        val eventSchema = new Schema.Parser().parse(schema)
        try {
          for (latency <- latencies) {
            val latencyRecord = new GenericData.Record(eventSchema)
            latencyRecord.put("eventname", latency.eventname)
            latencyRecord.put("second", latency.second)
            latencyRecord.put("operation", latency.operation)
            latencyRecord.put("value", latency.value)
            latencyRecord.put("ntpstatus", latency.ntpstatus)
            latencyRecord.put("cnt", latency.cnt)
            val record = new ProducerRecord[Any, AnyRef]("%s-latencies".format(topic), latencyRecord)
            producer.send(record)
          }
        } finally {
          producer.close()
        }
      })
    })

    latencyStream.foreachRDD(rdd => {
      rdd.saveToCassandra("spark_analysis", "events")
    })
  }
  
  def windowDuration(unit: String, durationValue: Long): Duration = unit match {
    case "second" => Seconds(durationValue)
    case "minute" => Minutes(durationValue)
  }
}

case class Event(eventname: String, second: Long, operation: String, value: Long, ntpstatus: Long, cnt: Long)
case class AppConfig(topic: String = "", brokerList: String = "", zookeeper: String = "", schemaRegistryUrl: String = "")
