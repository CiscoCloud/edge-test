/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ly.stealth.mesos.logging

import java.net.InetSocketAddress

import com.twitter.finagle.Service
import com.twitter.finagle.builder.{Server, ServerBuilder}
import com.twitter.finagle.http.{Http, Request, Response, RichHttp}
import com.twitter.util.Future
import org.apache.log4j.Logger
import org.apache.mesos.MesosExecutorDriver
import org.apache.mesos.Protos._

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

  val service: Service[Request, Response] = new Service[Request, Response] {
    def apply(req: Request): Future[Response] = {
      req.headerMap.get("Content-Type") match {
        case Some(contentType) =>
          if (config.sync) {
            new Thread {
              override def run() {
                transformer.transform(req.getContent().array(), contentType)
              }
            }.start()
          } else transformer.transform(req.getContent().array(), contentType)
        case None => logger.warn("no Content-Type header provided")
      }
      Future.value(Response())
    }
  }

  val server: Server = ServerBuilder()
    .codec(RichHttp[Request](Http()))
    .bindTo(new InetSocketAddress(config.port))
    .name("LogLine Transform")
    .build(service)
}

