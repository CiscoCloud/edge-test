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

import java.io.{ByteArrayOutputStream, InputStream}

import org.apache.log4j.Logger
import org.apache.mesos.MesosExecutorDriver
import org.apache.mesos.Protos._
import unfiltered.response.ResponseString

object Executor extends ExecutorBase {
  private var config: ExecutorConfigBase = null

  def main(args: Array[String]) {
    config = parseExecutorConfig(args)

    val driver = new MesosExecutorDriver(Executor)
    val status = if (driver.run eq Status.DRIVER_STOPPED) 0 else 1

    System.exit(status)
  }

  override protected def start() {
    new ExecutorEndpoint(config)
  }
}

class ExecutorEndpoint(config: ExecutorConfigBase) {
  private val logger = Logger.getLogger(this.getClass)
  private val transformer = new Transform(config)

  val plan = unfiltered.netty.cycle.Planify {
    case request =>
      request.headers("Content-Type").toList.headOption match {
        case Some(contentType) =>
          if (!config.sync) {
            new Thread {
              override def run() {
                transformer.transform(toBytes(request.inputStream), contentType, "Unfiltered")
              }
            }.start()
          } else transformer.transform(toBytes(request.inputStream), contentType, "Unfiltered")
        case None => logger.warn("no Content-Type header provided")
      }
      ResponseString("")
  }

  new Thread {
    override def run() {
      unfiltered.netty.Server.http(config.port).plan(plan).run()
    }
  }.start()

  private def toBytes(is: InputStream): Array[Byte] = {
    val buffer = new ByteArrayOutputStream()
    val data = new Array[Byte](16384)

    var nRead = -1
    while ( {
      nRead = is.read(data, 0, data.length)
      nRead != -1
    }) {
      buffer.write(data, 0, nRead)
    }

    buffer.flush()
    buffer.toByteArray
  }
}

