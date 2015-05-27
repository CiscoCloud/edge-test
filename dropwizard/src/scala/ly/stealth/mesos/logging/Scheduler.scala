package ly.stealth.mesos.logging

import java.util.UUID

import org.apache.mesos.Protos
import org.apache.mesos.Protos._

object Scheduler extends SchedulerBase {
  private var schedulerConfig: SchedulerConfig = null

  override protected def config: SchedulerConfigBase = schedulerConfig.base

  def parseConfig(args: Array[String]) {
    val parser = new scopt.OptionParser[SchedulerConfig]("scheduler") {
      override def errorOnUnknownArgument = false

      opt[String]('d', "executor.dropwizard.config").optional().text("Executor dropwizard config yml file.").action { (value, config) =>
        config.copy(executorDropwizardConfig = value)
      }
    }

    parser.parse(args, SchedulerConfig(base = parseSchedulerConfig(args))) match {
      case Some(c) => this.schedulerConfig = c
      case None => sys.exit(1)
    }
  }

  def main(args: Array[String]) {
    parseConfig(args)

    start(schedulerConfig.base, "Dropwizard LogLine Transform Framework")
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

    if (cpus > schedulerConfig.base.cpuPerTask && mems > schedulerConfig.base.memPerTask && portOpt.nonEmpty && adminPortOpt.nonEmpty) {
      val id = "transform-" + UUID.randomUUID().toString
      val taskId = TaskID.newBuilder().setValue(id).build()
      val taskInfo = TaskInfo.newBuilder().setName(taskId.getValue).setTaskId(taskId).setSlaveId(offer.getSlaveId)
        .setExecutor(this.createExecutor(id, portOpt.get, adminPortOpt.get))
        .addResources(Protos.Resource.newBuilder().setName("cpus").setType(Protos.Value.Type.SCALAR).setScalar(Protos.Value.Scalar.newBuilder().setValue(schedulerConfig.base.cpuPerTask)))
        .addResources(Protos.Resource.newBuilder().setName("mem").setType(Protos.Value.Type.SCALAR).setScalar(Protos.Value.Scalar.newBuilder().setValue(schedulerConfig.base.memPerTask)))
        .addResources(Protos.Resource.newBuilder().setName("ports").setType(Protos.Value.Type.RANGES).setRanges(
        Protos.Value.Ranges.newBuilder().addRange(Protos.Value.Range.newBuilder().setBegin(portOpt.get).setEnd(adminPortOpt.get))
      )).build

      Some(taskInfo)
    } else None
  }

  private def createExecutor(id: String, port: Long, adminPort: Long): ExecutorInfo = {
    val path = this.schedulerConfig.base.executor.split("/").last
    val executorConfigPath = this.schedulerConfig.executorDropwizardConfig.split("/").last
    val producerConfigPath = this.schedulerConfig.base.producerConfig.split("/").last
    val cmd = s"java -Ddw.server.applicationConnectors[0].port=$port -Ddw.server.adminConnectors[0].port=$adminPort -cp ${this.schedulerConfig.base.executor} ly.stealth.mesos.logging.Executor " +
      s"--producer.config ${this.schedulerConfig.base.producerConfig} --topic ${this.schedulerConfig.base.topic}"
    ExecutorInfo.newBuilder().setExecutorId(ExecutorID.newBuilder().setValue(id))
      .setCommand(CommandInfo.newBuilder()
      .addUris(CommandInfo.URI.newBuilder.setValue(s"http://${this.schedulerConfig.base.artifactServerHost}:${this.schedulerConfig.base.artifactServerPort}/resource/$path"))
      .addUris(CommandInfo.URI.newBuilder.setValue(s"http://${this.schedulerConfig.base.artifactServerHost}:${this.schedulerConfig.base.artifactServerPort}/resource/$executorConfigPath"))
      .addUris(CommandInfo.URI.newBuilder.setValue(s"http://${this.schedulerConfig.base.artifactServerHost}:${this.schedulerConfig.base.artifactServerPort}/resource/$producerConfigPath"))
      .setValue(cmd))
      .setName("LogLine Transform Executor")
      .setSource("cisco")
      .build
  }
}

private case class SchedulerConfig(base: SchedulerConfigBase = null, executorDropwizardConfig: String = "executor.yml")