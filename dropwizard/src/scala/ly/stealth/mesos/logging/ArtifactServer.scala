package ly.stealth.mesos.logging

import javax.ws.rs.core.Response
import javax.ws.rs.{GET, Path, PathParam}

import com.codahale.metrics.health.HealthCheck
import com.codahale.metrics.health.HealthCheck.Result
import io.dropwizard.setup.Environment
import io.dropwizard.{Application => App, Configuration}

import scala.io.Source

object ArtifactServer extends App[Configuration] {
  override def run(configuration: Configuration, environment: Environment) {
    environment.jersey().register(new Resource())
    environment.healthChecks().register("empty", EmptyHealthCheck)
  }
}

@Path("/")
class Resource {
  @GET
  @Path("{resource}")
  def getResource(@PathParam("resource") resource: String): Response = {
    val fileContents = Source.fromFile(resource).map(_.toByte).toArray
    Response.ok(fileContents).build
  }
}

object EmptyHealthCheck extends HealthCheck {
  override def check(): Result = Result.healthy()
}
