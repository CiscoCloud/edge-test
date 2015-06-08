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

import javax.ws.rs.{HeaderParam, POST, Path}

import com.codahale.metrics.health.HealthCheck
import com.codahale.metrics.health.HealthCheck.Result
import io.dropwizard.setup.Environment
import io.dropwizard.{Application => App, Configuration}

class ExecutorEndpoint(config: ExecutorConfig) extends App[Configuration] {

  override def run(configuration: Configuration, environment: Environment) {
    environment.jersey().register(new Handler(config))
    environment.healthChecks().register("empty", EmptyHealthCheck)
  }
}

@Path("/")
class Handler(config: ExecutorConfig) {
  private val transformer = new Transform(config.base)

  @POST
  def handle(body: Array[Byte], @HeaderParam("Content-Type") contentType: String) {
    if (!config.base.sync) {
      new Thread {
        override def run() {
          transformer.transform(body, contentType, "Dropwizard")
        }
      }.start()
    } else transformer.transform(body, contentType, "Dropwizard")
  }
}

object EmptyHealthCheck extends HealthCheck {
  override def check(): Result = Result.healthy()
}
