package ly.stealth.mesos.logging

import java.util

import ly.stealth.mesos.logging.Util.Str
import org.apache.log4j._
import org.apache.mesos.Protos._
import org.apache.mesos.{MesosSchedulerDriver, SchedulerDriver}

import scala.collection.JavaConversions._

object Scheduler extends org.apache.mesos.Scheduler {
  private val logger = Logger.getLogger(this.getClass)
  private var driver: SchedulerDriver = null

  private var runningInstances = 0
  private var config: SchedulerConfig = null

  def parseConfig(args: Array[String]) {
    val parser = new scopt.OptionParser[SchedulerConfig]("scheduler") {
      opt[String]('m', "master").required().text("Mesos Master addresses.").action {
        (value, config) =>
          config.copy(master = value)
      }

      opt[String]('u', "user").required().text("Mesos user.").action {
        (value, config) =>
          config.copy(user = value)
      }

      opt[Int]('i', "instances").required().text("Number of tasks to run.").action {
        (value, config) =>
          config.copy(instances = value)
      }

      opt[String]('h', "artifact.host").optional().text("Binding host for artifact server.").action {
        (value, config) =>
          config.copy(artifactServerHost = value)
      }

      opt[Int]('p', "artifact.port").optional().text("Binding port for artifact server.").action {
        (value, config) =>
          config.copy(artifactServerPort = value)
      }

      opt[String]('e', "executor").required().text("Executor file name.").action {
        (value, config) =>
          config.copy(executor = value)
      }

      opt[Double]('c', "cpu.per.task").optional().text("CPUs per task.").action {
        (value, config) =>
          config.copy(cpuPerTask = value)
      }

      opt[Double]('r', "mem.per.task").optional().text("Memory per task.").action {
        (value, config) =>
          config.copy(memPerTask = value)
      }

      opt[String]('s', "schema.registry").required().text("Avro Schema Registry url.").action {
        (value, config) =>
          config.copy(schemaRegistryUrl = value)
      }

      opt[String]('b', "broker.list").required().text("Comma separated list of brokers for producer.").action {
        (value, config) =>
          config.copy(brokerList = value)
      }

      opt[String]('t', "topic").required().text("Topic to produce transformed data to.").action {
        (value, config) =>
          config.copy(topic = value)
      }
    }

    parser.parse(args, SchedulerConfig()) match {
      case Some(c) => this.config = c
      case None => sys.exit(1)
    }
  }

  def main(args: Array[String]) {
    parseConfig(args)

    val frameworkBuilder = FrameworkInfo.newBuilder()
    frameworkBuilder.setUser(config.user)
    frameworkBuilder.setName("Dropwizard LogLine Transform Framework")

    val driver = new MesosSchedulerDriver(Scheduler, frameworkBuilder.build, config.master)

    Runtime.getRuntime.addShutdownHook(new Thread() {
      override def run() {
        if (driver != null) driver.stop()
      }
    })

    val status = if (driver.run eq Status.DRIVER_STOPPED) 0 else 1
    System.exit(status)
  }

  override def registered(driver: SchedulerDriver, id: FrameworkID, master: MasterInfo) {
    logger.info("[registered] framework:" + Str.id(id.getValue) + " master:" + Str.master(master))
    this.driver = driver
  }

  override def offerRescinded(driver: SchedulerDriver, id: OfferID) {
    logger.info("[offerRescinded] " + Str.id(id.getValue))
  }

  override def disconnected(driver: SchedulerDriver) {
    logger.info("[disconnected]")
    this.driver = null
  }

  override def reregistered(driver: SchedulerDriver, master: MasterInfo) {
    logger.info("[reregistered] master:" + Str.master(master))
    this.driver = driver
  }

  override def slaveLost(driver: SchedulerDriver, id: SlaveID) {
    logger.info("[slaveLost] " + Str.id(id.getValue))
  }

  override def error(driver: SchedulerDriver, message: String) {
    logger.info("[error] " + message)
  }

  override def statusUpdate(driver: SchedulerDriver, status: TaskStatus) {
    logger.info("[statusUpdate] " + Str.taskStatus(status))

    status.getState match {
      case TaskState.TASK_LOST | TaskState.TASK_FINISHED | TaskState.TASK_FAILED |
        TaskState.TASK_KILLED | TaskState.TASK_ERROR => this.runningInstances -= 1
      case _ =>
    }
  }

  override def frameworkMessage(driver: SchedulerDriver, executorId: ExecutorID, slaveId: SlaveID, data: Array[Byte]) {
    logger.info("[frameworkMessage] executor:" + Str.id(executorId.getValue) + " slave:" + Str.id(slaveId.getValue) + " data: " + new String(data))
  }

  override def resourceOffers(driver: SchedulerDriver, offers: util.List[Offer]) {
    logger.info("[resourceOffers]\n" + Str.offers(offers))
  }

  override def executorLost(driver: SchedulerDriver, executorId: ExecutorID, slaveId: SlaveID, status: Int) {
    logger.info("[executorLost] executor:" + Str.id(executorId.getValue) + " slave:" + Str.id(slaveId.getValue) + " status:" + status)
  }
}

private case class SchedulerConfig(master: String = "", user: String = "root", instances: Int = 1,
                                   artifactServerHost: String = "master", artifactServerPort: Int = 6666,
                                   executor: String = "", cpuPerTask: Double = 0.2, memPerTask: Double = 256,
                                   schemaRegistryUrl: String = "", brokerList: String = "", topic: String = "")
