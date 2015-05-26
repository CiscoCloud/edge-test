package ly.stealth.mesos.logging

import java.util

import ly.stealth.mesos.logging.Util.Str
import org.apache.log4j.Logger
import org.apache.mesos.Protos._
import org.apache.mesos.{Protos, Scheduler, SchedulerDriver}

import scala.collection.JavaConversions._
import scala.collection.mutable

abstract class SchedulerBase extends Scheduler {
  private val logger = Logger.getLogger(this.getClass)

  private var runningInstances = 0
  private val tasks: mutable.Set[TaskID] = mutable.Set()

  override def registered(driver: SchedulerDriver, id: FrameworkID, master: MasterInfo) {
    logger.info("[registered] framework:" + Str.id(id.getValue) + " master:" + Str.master(master))
  }

  override def offerRescinded(driver: SchedulerDriver, id: OfferID) {
    logger.info("[offerRescinded] " + Str.id(id.getValue))
  }

  override def disconnected(driver: SchedulerDriver) {
    logger.info("[disconnected]")
  }

  override def reregistered(driver: SchedulerDriver, master: MasterInfo) {
    logger.info("[reregistered] master:" + Str.master(master))
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
           TaskState.TASK_KILLED | TaskState.TASK_ERROR => synchronized {
        tasks -= status.getTaskId
        this.runningInstances -= 1
      }
      case _ =>
    }
  }

  override def frameworkMessage(driver: SchedulerDriver, executorId: ExecutorID, slaveId: SlaveID, data: Array[Byte]) {
    logger.info("[frameworkMessage] executor:" + Str.id(executorId.getValue) + " slave:" + Str.id(slaveId.getValue) + " data: " + new String(data))
  }

  override def resourceOffers(driver: SchedulerDriver, offers: util.List[Offer]) {
    logger.info("[resourceOffers]\n" + Str.offers(offers))

    synchronized {
      if (runningInstances > config.instances) {
        val toKill = runningInstances - config.instances
        ((0 until toKill) zip tasks).foreach { case (index, taskId) =>
          driver.killTask(taskId)
        }
      }
    }

    offers.foreach { offer =>
      synchronized {
        if (runningInstances < config.instances) {
          launchTask(offer) match {
            case Some(taskInfo) =>
              tasks += taskInfo.getTaskId
              runningInstances += 1
              driver.launchTasks(util.Arrays.asList(offer.getId), util.Arrays.asList(taskInfo), Filters.newBuilder().setRefuseSeconds(1).build)
            case None => driver.declineOffer(offer.getId)
          }
        }
      }
    }
  }

  override def executorLost(driver: SchedulerDriver, executorId: ExecutorID, slaveId: SlaveID, status: Int) {
    logger.info("[executorLost] executor:" + Str.id(executorId.getValue) + " slave:" + Str.id(slaveId.getValue) + " status:" + status)
  }

  protected def getScalarResources(offer: Offer, name: String): Double = {
    offer.getResourcesList.foldLeft(0.0) { (all, current) =>
      if (current.getName == name) all + current.getScalar.getValue
      else all
    }
  }

  protected def getRangeResources(offer: Offer, name: String): List[Protos.Value.Range] = {
    offer.getResourcesList.foldLeft[List[Protos.Value.Range]](List()) { case (all, current) =>
      if (current.getName == name) all ++ current.getRanges.getRangeList
      else all
    }
  }

  protected def config: SchedulerConfigBase

  protected def launchTask(offer: Offer): Option[TaskInfo]
}

trait SchedulerConfigBase {
  var cpuPerTask: Double = 0.2
  var memPerTask: Double = 256
  var instances: Int = 1
}
