package ly.stealth.mesos.logging

import java.util.UUID

import org.apache.mesos.Protos._
import org.apache.mesos.{MesosSchedulerDriver, Protos}

object Scheduler extends SchedulerBase {
  private var schedulerConfig: SchedulerConfig = null

  override protected def config: SchedulerConfigBase = schedulerConfig

  def parseConfig(args: Array[String]) {
    val parser = new scopt.OptionParser[SchedulerConfig]("scheduler") {
      opt[String]('m', "master").required().text("Mesos Master addresses.").action { (value, config) =>
        config.copy(master = value)
      }

      opt[String]('u', "user").required().text("Mesos user.").action { (value, config) =>
        config.copy(user = value)
      }

      opt[Int]('i', "instances").optional().text("Number of tasks to run.").action { (value, config) =>
        config.instances = value
        config
      }

      opt[String]('h', "artifact.host").optional().text("Binding host for artifact server.").action { (value, config) =>
        config.copy(artifactServerHost = value)
      }

      opt[Int]('p', "artifact.port").optional().text("Binding port for artifact server.").action { (value, config) =>
        config.copy(artifactServerPort = value)
      }

      opt[String]('e', "executor").required().text("Executor file name.").action { (value, config) =>
        config.copy(executor = value)
      }

      opt[Double]('c', "cpu.per.task").optional().text("CPUs per task.").action { (value, config) =>
        config.cpuPerTask = value
        config
      }

      opt[Double]('r', "mem.per.task").optional().text("Memory per task.").action { (value, config) =>
        config.memPerTask = value
        config
      }

      opt[String]('s', "producer.config").required().text("Producer config file name.").action { (value, config) =>
        config.copy(producerConfig = value)
      }

      opt[String]('t', "topic").required().text("Topic to produce transformed data to.").action { (value, config) =>
        config.copy(topic = value)
      }

      opt[String]('d', "executor.dropwizard.config").optional().text("Executor dropwizard config yml file.").action { (value, config) =>
        config.copy(executorDropwizardConfig = value)
      }
    }

    parser.parse(args, SchedulerConfig()) match {
      case Some(c) => this.schedulerConfig = c
      case None => sys.exit(1)
    }
  }

  def main(args: Array[String]) {
    parseConfig(args)

    val server = new HttpServer(schedulerConfig.artifactServerPort, schedulerConfig)

    val frameworkBuilder = FrameworkInfo.newBuilder()
    frameworkBuilder.setUser(schedulerConfig.user)
    frameworkBuilder.setName("Dropwizard LogLine Transform Framework")

    val driver = new MesosSchedulerDriver(Scheduler, frameworkBuilder.build, schedulerConfig.master)

    Runtime.getRuntime.addShutdownHook(new Thread() {
      override def run() {
        if (driver != null) driver.stop()
      }
    })

    val status = if (driver.run eq Status.DRIVER_STOPPED) 0 else 1
    server.stop()
    System.exit(status)
  }

  override def launchTask(offer: Offer): Option[TaskInfo] = {
    val cpus = getScalarResources(offer, "cpus")
    val mems = getScalarResources(offer, "mem")
    val ports = getRangeResources(offer, "ports")
    val portOpt = ports.headOption.map(_.getBegin)
    val adminPortOpt = ports.headOption.flatMap { range =>
      val port = range.getBegin + 1
      if (range.getEnd >= port) Some(port)
      else None
    }

    if (cpus > schedulerConfig.cpuPerTask && mems > schedulerConfig.memPerTask && portOpt.nonEmpty && adminPortOpt.nonEmpty) {
      val id = "transform-" + UUID.randomUUID().toString
      val taskId = TaskID.newBuilder().setValue(id).build()
      val taskInfo = TaskInfo.newBuilder().setName(taskId.getValue).setTaskId(taskId).setSlaveId(offer.getSlaveId)
        .setExecutor(this.createExecutor(id, portOpt.get, adminPortOpt.get))
        .addResources(Protos.Resource.newBuilder().setName("cpus").setType(Protos.Value.Type.SCALAR).setScalar(Protos.Value.Scalar.newBuilder().setValue(schedulerConfig.cpuPerTask)))
        .addResources(Protos.Resource.newBuilder().setName("mem").setType(Protos.Value.Type.SCALAR).setScalar(Protos.Value.Scalar.newBuilder().setValue(schedulerConfig.memPerTask)))
        .addResources(Protos.Resource.newBuilder().setName("ports").setType(Protos.Value.Type.RANGES).setRanges(
        Protos.Value.Ranges.newBuilder().addRange(Protos.Value.Range.newBuilder().setBegin(portOpt.get).setEnd(adminPortOpt.get))
      )).build

      Some(taskInfo)
    } else None
  }

  private def createExecutor(id: String, port: Long, adminPort: Long): ExecutorInfo = {
    val path = this.schedulerConfig.executor.split("/").last
    val executorConfigPath = this.schedulerConfig.executorDropwizardConfig.split("/").last
    val producerConfigPath = this.schedulerConfig.producerConfig.split("/").last
    val cmd = s"java -Ddw.server.applicationConnectors[0].port=$port -Ddw.server.adminConnectors[0].port=$adminPort -cp ${this.schedulerConfig.executor} ly.stealth.mesos.logging.Executor " +
      s"--producer.config ${this.schedulerConfig.producerConfig} --topic ${this.schedulerConfig.topic}"
    ExecutorInfo.newBuilder().setExecutorId(ExecutorID.newBuilder().setValue(id))
      .setCommand(CommandInfo.newBuilder()
      .addUris(CommandInfo.URI.newBuilder.setValue(s"http://${this.schedulerConfig.artifactServerHost}:${this.schedulerConfig.artifactServerPort}/resource/$path"))
      .addUris(CommandInfo.URI.newBuilder.setValue(s"http://${this.schedulerConfig.artifactServerHost}:${this.schedulerConfig.artifactServerPort}/resource/$executorConfigPath"))
      .addUris(CommandInfo.URI.newBuilder.setValue(s"http://${this.schedulerConfig.artifactServerHost}:${this.schedulerConfig.artifactServerPort}/resource/$producerConfigPath"))
      .setValue(cmd))
      .setName("LogLine Transform Executor")
      .setSource("cisco")
      .build
  }
}

private case class SchedulerConfig(master: String = "", user: String = "root", artifactServerHost: String = "master", artifactServerPort: Int = 6666,
                                   executor: String = "", producerConfig: String = "", topic: String = "",
                                   executorDropwizardConfig: String = "executor.yml") extends SchedulerConfigBase