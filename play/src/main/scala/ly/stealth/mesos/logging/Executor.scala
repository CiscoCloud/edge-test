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

import org.apache.log4j.Logger
import org.apache.mesos.MesosExecutorDriver
import org.apache.mesos.Protos._
import play.api.mvc._
import play.api.routing.sird._
import play.core.server._

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
  Thread.currentThread().setContextClassLoader(this.getClass.getClassLoader)

  private val logger = Logger.getLogger(this.getClass)
  private val transformer = new Transform(config)

  NettyServer.fromRouter(ServerConfig(
    port = Some(config.port)
  )) {
    case POST(p"/") => Action(BodyParsers.parse.raw) { request =>
      val data = for {
        contentType <- request.headers.get("Content-Type")
        body <- request.body.asBytes()
      } yield (contentType, body)

      data match {
        case Some((contentType, body)) => transformer.transform(body, contentType)
        case None => logger.warn("Either no Content-Type header provided or body is empty")
      }

      Results.Ok
    }
  }
}

