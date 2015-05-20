package ly.stealth.mesos.logging

import java.io.File
import javax.ws.rs.core.Response
import javax.ws.rs.{GET, Path, PathParam}

import com.codahale.metrics.health.HealthCheck
import com.codahale.metrics.health.HealthCheck.Result
import io.dropwizard.setup.Environment
import io.dropwizard.{Application => App, Configuration}

class HttpServer(config: SchedulerConfig) extends App[Configuration] {
  override def run(configuration: Configuration, environment: Environment) {
    environment.jersey().register(new Resource())
    environment.jersey().register(new Scale(config))
    environment.healthChecks().register("empty", EmptyHealthCheck)
  }
}

@Path("/resource/")
class Resource {
  @GET
  @Path("{resource}")
  def getResource(@PathParam("resource") resource: String): Response = {
    Response.ok(new File(resource)).build
  }
}

@Path("/scale/")
class Scale(config: SchedulerConfig) {
  @GET
  @Path("{scale}")
  def scale(@PathParam("scale") scale: Int): Response = {
    if (scale >= 0) config.instances = scale
    Response.ok().build
  }
}

object EmptyHealthCheck extends HealthCheck {
  override def check(): Result = Result.healthy()
}
