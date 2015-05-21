package ly.stealth.mesos.logging

import ly.stealth.mesos.logging.Util.Str
import org.apache.log4j.Logger
import org.apache.mesos.Protos._
import org.apache.mesos.{ExecutorDriver, MesosExecutorDriver, Protos}

object Executor extends org.apache.mesos.Executor {
  private var shutdownFlag = false
  private var config: ExecutorConfig = null
  private val lock = new Object()

  def parseConfig(args: Array[String]) {
    val parser = new scopt.OptionParser[ExecutorConfig]("executor") {
      opt[String]('p', "producer.config").required().text("Producer config file name.").action { (value, config) =>
        config.copy(producerConfig = value)
      }

      opt[String]('t', "topic").required().text("Topic to produce transformed data to.").action { (value, config) =>
        config.copy(topic = value)
      }

      opt[String]('d', "dropwizard.config").optional().text("Dropwizard config yml file.").action { (value, config) =>
        config.copy(dropwizardConfig = value)
      }
    }

    parser.parse(args, ExecutorConfig()) match {
      case Some(c) => this.config = c
      case None => sys.exit(1)
    }
  }

  def main(args: Array[String]) {
    parseConfig(args)

    val driver = new MesosExecutorDriver(Executor)
    val status = if (driver.run eq Status.DRIVER_STOPPED) 0 else 1

    System.exit(status)
  }

  private val logger = Logger.getLogger(this.getClass)

  override def registered(driver: ExecutorDriver, executorInfo: ExecutorInfo, framework: FrameworkInfo, slave: SlaveInfo) {
    logger.info("[registered] framework:" + Str.framework(framework) + " slave:" + Str.slave(slave))
  }

  override def shutdown(driver: ExecutorDriver) {
    logger.info("[shutdown]")

    lock.synchronized {
      shutdownFlag = true
      lock.notifyAll()
    }
  }

  override def disconnected(driver: ExecutorDriver) {
    logger.info("[disconnected]")
  }

  override def killTask(driver: ExecutorDriver, id: TaskID) {
    logger.info("[killTask] " + id.getValue)

    lock.synchronized {
      shutdownFlag = true
      lock.notifyAll()
    }
  }

  override def reregistered(driver: ExecutorDriver, slave: SlaveInfo) {
    logger.info("[reregistered] " + Str.slave(slave))
  }

  override def error(driver: ExecutorDriver, message: String) {
    logger.info("[error] " + message)
  }

  override def frameworkMessage(driver: ExecutorDriver, data: Array[Byte]) {
    logger.info("[frameworkMessage] " + new String(data))
  }

  override def launchTask(driver: ExecutorDriver, task: TaskInfo) {
    logger.info("[launchTask] " + Str.task(task))

    val runStatus = TaskStatus.newBuilder().setTaskId(task.getTaskId).setState(Protos.TaskState.TASK_RUNNING).build

    driver.sendStatusUpdate(runStatus)

    new Thread {
      override def run() {
        var failed = false
        try {
          logger.info("[Started task] " + Str.task(task))
          new ExecutorEndpoint(config).run("server", config.dropwizardConfig)

          while (!shutdownFlag) {
            lock.synchronized(lock.wait())
          }
        } catch {
          case t: Throwable =>
            logger.warn("", t)
            failed = true
        } finally {
          val finishedStatus = TaskStatus.newBuilder().setTaskId(task.getTaskId)
            .setState(if (failed) Protos.TaskState.TASK_FAILED else Protos.TaskState.TASK_FINISHED).build
          driver.sendStatusUpdate(finishedStatus)
          driver.stop()
        }
      }
    }.start()
  }
}

case class ExecutorConfig(producerConfig: String = "", topic: String = "", dropwizardConfig: String = "executor.yml")


