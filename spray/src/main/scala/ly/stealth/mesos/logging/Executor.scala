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
    case HttpRequest(POST, Uri.Path("/"), headers, entity: HttpEntity.NonEmpty, _) => {
      headers.find(_.is("content-type")).foreach { contentType =>
        transformer.transform(entity.data.toByteArray, contentType.value)
      }
      sender ! HttpResponse()
    }
  }
}

