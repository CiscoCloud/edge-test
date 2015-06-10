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
import org.apache.spark.streaming.kafka._

object Main extends App {
  val parser = new scopt.OptionParser[AppConfig]("spark-analysis") {
    head("Latencies calculation job", "1.0")
    opt[String]("topics") unbounded() required() action { (value, config) =>
      config.copy(topic = value)
    } text ("Comma separated list of topics to read data from")
    opt[String]("zookeeper") unbounded() required() action { (value, config) =>
      config.copy(zookeeper = value)
    } text ("Zookeeper connection string - host:port")
    opt[String]("broker.list") unbounded() required() action { (value, config) =>
      config.copy(brokerList = value)
    } text ("Comma separated string of host:port")
    opt[String]("schema.registry.url") unbounded() required() action { (value, config) =>
      config.copy(schemaRegistryUrl = value)
    } text ("Schema registry URL")
    opt[Int]("partitions") unbounded() required() action { (value, config) =>
      config.copy(partitions = value)
    } text ("Initial amount of RDD partitions")
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
    session.execute("CREATE TABLE IF NOT EXISTS spark_analysis.events(framework text, second bigint, message_size bigint, eventname text, latency double, received_count int, sent_count int, PRIMARY KEY(framework, second, message_size, eventname)) WITH CLUSTERING ORDER BY (second DESC)")
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

  start(ssc, consumerConfig, producerConfig, appConfig.topic, appConfig.partitions)

  ssc.start()
  ssc.awaitTermination()

  def start(ssc: StreamingContext, consumerConfig: Map[String, String], producerConfig: Properties, topics: String, partitions: Int) = {
    val topicMap = topics.split(",").map(_ -> partitions).toMap
    val latencyStream = KafkaUtils.createStream[Array[Byte], SchemaAndData, DefaultDecoder, AvroDecoder](ssc, consumerConfig, topicMap, StorageLevel.MEMORY_ONLY).map(value => {
      val record = value._2.deserialize().asInstanceOf[GenericRecord]
      import scala.collection.JavaConversions._
      val timings = record.get("timings").asInstanceOf[GenericData.Array[Record]]
      val topic = record.get("tag").asInstanceOf[java.util.Map[Utf8, Utf8]].get(new Utf8("topic")).toString
      (timings.head.get("eventName").asInstanceOf[Utf8].toString + "-" + timings.last.get("eventName").asInstanceOf[Utf8].toString,
       timings.head.get("value").asInstanceOf[Long] / 1000000000,
       timings.last.get("value").asInstanceOf[Long] / 1000000000,
       (timings.last.get("value").asInstanceOf[Long] - timings.head.get("value").asInstanceOf[Long]).toDouble / 1000000,
       record.get("source").asInstanceOf[Utf8].toString,
       record.get("size").asInstanceOf[Long],
       topic)
    }).transform( rdd => {
      rdd.groupBy(entry => (entry._1, entry._3, entry._5, entry._6, entry._7))
    }.map( entry => {
      val key = entry._1
      val values = entry._2
      val receivedValuesCount = values.count(item => item._2 == item._3).toLong
      val avgLatency = values.map(_._4).sum / values.size
      (key._3, key._2, key._4, key._1, avgLatency, receivedValuesCount, values.size.toLong, key._5)
    })).persist()

    val schema = "{\"type\":\"record\",\"name\":\"event\",\"fields\":[{\"name\":\"framework\",\"type\":\"string\"},{\"name\":\"second\",\"type\":\"long\"},{\"name\":\"message_size\",\"type\":\"long\"},{\"name\":\"eventname\",\"type\":\"string\"},{\"name\":\"latency\",\"type\":\"double\"},{\"name\":\"received_count\",\"type\":\"long\"},{\"name\":\"sent_count\",\"type\":\"long\"}]}"
    latencyStream.foreachRDD(rdd => {
      rdd.foreachPartition(events => {
        val producer = new KafkaProducer[Any, AnyRef](producerConfig)
        val eventSchema = new Schema.Parser().parse(schema)
        try {
          for (event <- events) {
            val latencyRecord = new GenericData.Record(eventSchema)
            latencyRecord.put("framework", event._1)
            latencyRecord.put("second", event._2)
            latencyRecord.put("message_size", event._3)
            latencyRecord.put("eventname", event._4)
            latencyRecord.put("latency", event._5)
            latencyRecord.put("received_count", event._6)
            latencyRecord.put("sent_count", event._7)
            val record = new ProducerRecord[Any, AnyRef]("%s-latencies".format(event._8), latencyRecord)
            producer.send(record).get()
          }
        } finally {
          producer.close()
        }
      })
    })

    latencyStream.foreachRDD(rdd => {
      rdd.saveToCassandra("spark_analysis", "events", SomeColumns("framework", "second", "message_size", "eventname", "latency", "received_count", "sent_count"))
    })
  }
}

case class AppConfig(topic: String = "", brokerList: String = "", zookeeper: String = "", partitions: Int = 1, schemaRegistryUrl: String = "")
