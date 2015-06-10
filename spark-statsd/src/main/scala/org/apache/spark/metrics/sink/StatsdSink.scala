package org.apache.spark.metrics.sink

import java.util.Properties
import java.util.concurrent.TimeUnit

import com.bealetech.metrics.reporting.{Statsd, StatsdReporter}
import com.codahale.metrics.MetricRegistry
import org.apache.spark.metrics.MetricsSystem

private[spark] class StatsdSink(val property: Properties, val registry: MetricRegistry) extends Sink {
  val STATSD_DEFAULT_PERIOD = 10
  val STATSD_DEFAULT_RATE_UNIT = "SECONDS"
  val STATSD_DEFAULT_DURATION_UNIT = "MILLISECONDS"
  val STATSD_DEFAULT_PREFIX = ""
  
  val STATSD_KEY_HOST = "host"
  val STATSD_KEY_PORT = "port"
  val STATSD_KEY_PERIOD = "period"
  val STATSD_KEY_RATE_UNIT = "rateUnit"
  val STATSD_KEY_DURATION_UNIT = "durationUnit"
  val STATSD_KEY_PREFIX = "prefix"

  def propertyToOption(prop: String): Option[String] = Option(property.getProperty(prop))

  if (!propertyToOption(STATSD_KEY_HOST).isDefined) {
    throw new Exception("Statsd sink requires 'host' property.")
  }

  if (!propertyToOption(STATSD_KEY_PORT).isDefined) {
    throw new Exception("Statsd sink requires 'port' property.")
  }

  val host = propertyToOption(STATSD_KEY_HOST).get
  val port = propertyToOption(STATSD_KEY_PORT).get.toInt

  val rate = propertyToOption(STATSD_KEY_PERIOD) match {
    case Some(s) => s.toInt
    case None => STATSD_DEFAULT_PERIOD
  }

  val rateUnit: TimeUnit = propertyToOption(STATSD_KEY_RATE_UNIT) match {
    case Some(s) => TimeUnit.valueOf(s.toUpperCase)
    case None => TimeUnit.valueOf(STATSD_DEFAULT_RATE_UNIT)
  }

  val durationUnit: TimeUnit = propertyToOption(STATSD_KEY_DURATION_UNIT) match {
    case Some(s) => TimeUnit.valueOf(s.toUpperCase)
    case None => TimeUnit.valueOf(STATSD_DEFAULT_DURATION_UNIT)
  }

  val prefix = propertyToOption(STATSD_KEY_PREFIX).getOrElse(STATSD_DEFAULT_PREFIX)

  MetricsSystem.checkMinimalPollingPeriod(rateUnit, rate)

  val reporter = StatsdReporter.forRegistry(registry)
    .convertDurationsTo(durationUnit)
    .convertRatesTo(rateUnit)
    .prefixedWith(prefix).build(new Statsd(host, port))

  def start {
    reporter.start(rate, rateUnit)
  }

  def stop {
    reporter.stop()
  }

  def report() {
    reporter.report()
  }
}
