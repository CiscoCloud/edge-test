package ly.stealth.mesos.logging

import ly.stealth.mesos.logging.Util.Str
import org.apache.log4j.Logger
import org.apache.mesos.{MesosExecutorDriver, Protos, ExecutorDriver}
import org.apache.mesos.Protos._

object Executor extends org.apache.mesos.Executor {
  def main(args: Array[String]) {
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
  }

  override def disconnected(driver: ExecutorDriver) {
    logger.info("[disconnected]")
  }

  override def killTask(driver: ExecutorDriver, id: TaskID) {
    logger.info("[killTask] " + id.getValue)
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
        logger.info("[Started task] " +  Str.task(task))
        Thread.sleep(30000)

        val finishedStatus = TaskStatus.newBuilder().setTaskId(task.getTaskId).setState(Protos.TaskState.TASK_FINISHED).build
        driver.sendStatusUpdate(finishedStatus)
      }
    }.start()
  }
}
