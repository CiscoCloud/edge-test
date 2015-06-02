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

import java.util

import ly.stealth.mesos.logging.Util.Str
import org.apache.log4j.Logger
import org.apache.mesos.Protos._
import org.apache.mesos.{MesosSchedulerDriver, Protos, Scheduler, SchedulerDriver}

import scala.collection.JavaConversions._
import scala.collection.mutable

abstract class SchedulerBase extends Scheduler {
  private val logger = Logger.getLogger(this.getClass)

  private var runningInstances = 0
  private val tasks: mutable.Set[TaskID] = mutable.Set()

  def parseSchedulerConfig(args: Array[String]): SchedulerConfigBase = {
    val parser = new scopt.OptionParser[SchedulerConfigBase]("scheduler") {
      override def errorOnUnknownArgument = false

      opt[String]('m', "master").required().text("Mesos Master addresses.").action { (value, config) =>
        config.copy(master = value)
      }

      opt[String]('u', "user").required().text("Mesos user.").action { (value, config) =>
        config.copy(user = value)
      }

      opt[Int]('i', "instances").optional().text("Number of tasks to run.").action { (value, config) =>
        config.copy(instances = value)
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
        config.copy(cpuPerTask = value)
      }

      opt[Double]('r', "mem.per.task").optional().text("Memory per task.").action { (value, config) =>
        config.copy(memPerTask = value)
      }

      opt[String]('s', "producer.config").required().text("Producer config file name.").action { (value, config) =>
        config.copy(producerConfig = value)
      }

      opt[String]('t', "topic").required().text("Topic to produce transformed data to.").action { (value, config) =>
        config.copy(topic = value)
      }

      opt[Boolean]('s', "sync").optional().text("Flag to respond only after decoding-encoding is done.").action { (value, config) =>
        config.copy(sync = value)
      }
    }

    parser.parse(args, SchedulerConfigBase()) match {
      case Some(config) => config
      case None => sys.exit(1)
    }
  }

  def start(config: SchedulerConfigBase, name: String) {
    val server = new HttpServer(config)

    val frameworkBuilder = FrameworkInfo.newBuilder()
    frameworkBuilder.setUser(config.user)
    frameworkBuilder.setName(name)

    val driver = new MesosSchedulerDriver(this, frameworkBuilder.build, config.master)

    Runtime.getRuntime.addShutdownHook(new Thread() {
      override def run() {
        if (driver != null) driver.stop()
      }
    })

    val status = if (driver.run eq Status.DRIVER_STOPPED) 0 else 1
    server.stop()
    System.exit(status)
  }

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
        } else driver.declineOffer(offer.getId)
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

case class SchedulerConfigBase(master: String = "", user: String = "root", cpuPerTask: Double = 0.2,
                               memPerTask: Double = 256, var instances: Int = 1, artifactServerHost: String = "master",
                               artifactServerPort: Int = 6666, executor: String = "", producerConfig: String = "",
                               topic: String = "", sync: Boolean = false)
