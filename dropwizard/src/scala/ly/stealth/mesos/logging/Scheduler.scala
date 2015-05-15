package ly.stealth.mesos.logging

import java.util

import org.apache.log4j._
import org.apache.mesos.Protos._
import org.apache.mesos.{MesosSchedulerDriver, SchedulerDriver}
import Util.Str
import scala.collection.JavaConversions._

object Scheduler extends org.apache.mesos.Scheduler {
  private val logger = Logger.getLogger(this.getClass)
  private var driver: SchedulerDriver = null

  def main(args: Array[String]) {
    val frameworkBuilder = FrameworkInfo.newBuilder()
    frameworkBuilder.setUser(args(1)) //TODO use flags
    frameworkBuilder.setName("Dropwizard LogLine Transform Framework")

    val driver = new MesosSchedulerDriver(Scheduler, frameworkBuilder.build, args(0)) //TODO use flags

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
