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
  private val transformer = new Transform(config)

  @POST
  def handle(body: Array[Byte], @HeaderParam("Content-Type") contentType: String) {
    transformer.transform(body, contentType)
  }
}

object EmptyHealthCheck extends HealthCheck {
  override def check(): Result = Result.healthy()
}
