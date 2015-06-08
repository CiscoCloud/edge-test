/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ly.stealth.mesos.logging

import java.io.FileInputStream
import java.util.Properties
import java.util.concurrent.TimeUnit

import _root_.io.confluent.kafka.serializers.{KafkaAvroDecoder, KafkaAvroSerializer}
import com.codahale.metrics.{ConsoleReporter, MetricRegistry}
import kafka.utils.VerifiableProperties
import org.apache.avro.generic.{GenericRecord, IndexedRecord}
import org.apache.kafka.clients.producer.ProducerConfig._
import org.apache.kafka.clients.producer.{KafkaProducer, ProducerRecord}
import org.apache.log4j.Logger
import org.codehaus.jackson.Version
import org.codehaus.jackson.map.module.SimpleModule
import org.codehaus.jackson.map.{DeserializationContext, KeyDeserializer, ObjectMapper}

import scala.collection.JavaConversions._
import scala.util.{Failure, Success, Try}

class Transform(config: ExecutorConfigBase) {
  private val metrics = new MetricRegistry
  private val requestsPerSec = metrics.meter("requests")
  
  ConsoleReporter.forRegistry(metrics).convertRatesTo(TimeUnit.SECONDS)
  .convertDurationsTo(TimeUnit.MILLISECONDS).build().start(1, TimeUnit.SECONDS)

  final val CONTENT_TYPE_AVRO = "avro/binary"
  final val CONTENT_TYPE_PROTOBUF = "application/x-protobuf"
  final val CONTENT_TYPE_JSON = "application/json"

  private val logger = Logger.getLogger(this.getClass)
  private val mapper = new ObjectMapper()
  private val module = new SimpleModule("charsequence-module", Version.unknownVersion())
  module.addKeyDeserializer(classOf[CharSequence], new CharSequenceKeyDeserializer)
  mapper.registerModule(module)

  private val props = new Properties()
  props.load(new FileInputStream(config.producerConfig))
  props.put(KEY_SERIALIZER_CLASS_CONFIG, classOf[KafkaAvroSerializer])
  props.put(VALUE_SERIALIZER_CLASS_CONFIG, classOf[KafkaAvroSerializer])
  logger.info("Producer properties: " + props)

  private val producer = new KafkaProducer[Any, IndexedRecord](props)

  private val avroDecoder = new KafkaAvroDecoder(new VerifiableProperties(props))

  def transform(data: Array[Byte], contentType: String, framework: String) {
    val received = timing("received")
    requestsPerSec.mark()

    val logLineOpt = contentType match {
      case CONTENT_TYPE_JSON => this.handleJson(data)
      case CONTENT_TYPE_PROTOBUF => this.handleProtobuf(data)
      case CONTENT_TYPE_AVRO => this.handleAvro(data)
      case _ =>
        logger.warn(s"Content-Type $contentType is invalid")
        None
    }

    logLineOpt.foreach { logLine =>
      logLine.setSize(data.length.toLong)
      logLine.setSource(framework)
      if (logLine.getTag == null) logLine.setTag(new java.util.HashMap[CharSequence, CharSequence])
      logLine.getTag.put("topic", config.topic)
      logLine.getTimings.add(received)
      logLine.getTimings.add(timing("sent"))

      producer.send(new ProducerRecord[Any, IndexedRecord](config.topic, logLine))
    }
  }

  private def handleJson(body: Array[Byte]): Option[LogLine] = {
    Try(mapper.readValue(body, classOf[LogLine])) match {
      case Success(logLine) => Some(logLine)
      case Failure(ex) =>
        logger.warn("", ex)
        None
    }
  }

  private def handleProtobuf(body: Array[Byte]): Option[LogLine] = {
    Try(proto.Logline.LogLine.parseFrom(body)) match {
      case Success(protoLine) =>
        val logLine = new LogLine()
        logLine.setLine(protoLine.getLine)
        logLine.setLogtypeid(protoLine.getLogtypeid)
        logLine.setSource(protoLine.getSource)
        logLine.setTag(mapAsJavaMap(protoLine.getTagList.map(tag => tag.getKey -> tag.getValue).toMap))
        logLine.setTimings(protoLine.getTimingsList.map(protoTiming => Timing.newBuilder().setEventName(protoTiming.getEventName).setValue(protoTiming.getValue).build))
        Some(logLine)
      case Failure(ex) =>
        logger.warn("", ex)
        None
    }
  }

  private def handleAvro(body: Array[Byte]): Option[LogLine] = {
    Try(avroDecoder.fromBytes(body)) match {
      case Success(obj) =>
        val generic = obj.asInstanceOf[GenericRecord]
        val logLine = new LogLine()
        logLine.setLine(generic.get("line").asInstanceOf[CharSequence])
        logLine.setLogtypeid(generic.get("logtypeid").asInstanceOf[java.lang.Long])
        logLine.setSource(generic.get("source").asInstanceOf[CharSequence])
        val tags = generic.get("tag")
        if (tags != null) logLine.setTag(tags.asInstanceOf[Map[CharSequence, CharSequence]])
        logLine.setTimings(generic.get("timings").asInstanceOf[java.util.List[GenericRecord]].map { timing =>
          Timing.newBuilder().setEventName(timing.get("eventName").asInstanceOf[CharSequence]).setValue(timing.get("value").asInstanceOf[java.lang.Long]).build
        })
        Some(logLine)
      case Failure(ex) =>
        logger.warn("", ex)
        None
    }
  }

  //TODO ntpstatus
  private def timing(name: String): Timing = Timing.newBuilder().setEventName(name)
    .setValue(System.currentTimeMillis() * 1000000 + System.nanoTime() % 1000000).build
}

class CharSequenceKeyDeserializer extends KeyDeserializer {
  override def deserializeKey(key: String, ctxt: DeserializationContext): AnyRef = key
}
