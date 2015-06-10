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

import java.net.InetAddress

import akka.actor.{Actor, ActorLogging, ActorSystem, Props}
import akka.io.IO
import org.apache.mesos.MesosExecutorDriver
import org.apache.mesos.Protos._
import spray.can.Http
import spray.http.HttpMethods._
import spray.http.{HttpRequest, _}

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
  implicit val system = ActorSystem("spray-transform")
  val service = system.actorOf(Props(new TransformActor(config)), "spray-transform-service")
  IO(Http) ! Http.Bind(service, interface = InetAddress.getLocalHost.getHostName, port = config.port)
}

class TransformActor(config: ExecutorConfigBase) extends Actor with ActorLogging {
  private val transformer = new Transform(config)

  def receive = {
    case _: Http.Connected => sender ! Http.Register(self)
    case HttpRequest(POST, Uri.Path("/"), headers, entity: HttpEntity.NonEmpty, _) =>
      headers.find(_.is("content-type")).foreach { contentType =>
        if (!config.sync) {
          new Thread {
            override def run() {
              transformer.transform(entity.data.toByteArray, contentType.value, "Spray")
            }
          }.start()
        } else transformer.transform(entity.data.toByteArray, contentType.value, "Spray")
      }
      sender ! HttpResponse()
  }
}

