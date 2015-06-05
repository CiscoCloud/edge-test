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

import org.apache.mesos.MesosExecutorDriver
import org.apache.mesos.Protos._

object Executor extends ExecutorBase {
  private var config: ExecutorConfig = null

  def parseConfig(args: Array[String]) {
    val parser = new scopt.OptionParser[ExecutorConfig]("executor") {
      override def errorOnUnknownArgument = false

      opt[String]('d', "dropwizard.config").optional().text("Dropwizard config yml file.").action { (value, config) =>
        config.copy(dropwizardConfig = value)
      }
    }

    parser.parse(args, ExecutorConfig(base = parseExecutorConfig(args))) match {
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

  override protected def name(): String = "Dropwizard"
}

case class ExecutorConfig(base: ExecutorConfigBase, dropwizardConfig: String = "executor.yml")


