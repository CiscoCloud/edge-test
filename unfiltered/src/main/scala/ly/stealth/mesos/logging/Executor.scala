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
        case Some(contentType) => transformer.transform(toBytes(request.inputStream), contentType)
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

