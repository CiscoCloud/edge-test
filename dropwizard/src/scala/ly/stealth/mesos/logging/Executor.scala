package ly.stealth.mesos.logging

import org.apache.mesos.MesosExecutorDriver
import org.apache.mesos.Protos._

object Executor extends ExecutorBase {
  private var config: ExecutorConfig = null

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

  override protected def start() {
    new ExecutorEndpoint(config).run("server", config.dropwizardConfig)
  }
}

case class ExecutorConfig(producerConfig: String = "", topic: String = "", dropwizardConfig: String = "executor.yml")


