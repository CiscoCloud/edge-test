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
        case Some(contentType) => transformer.transform(req.getContent().array(), contentType)
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

