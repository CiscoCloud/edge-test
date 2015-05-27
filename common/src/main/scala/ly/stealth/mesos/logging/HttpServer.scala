package ly.stealth.mesos.logging

import java.io._
import javax.servlet.http.{HttpServlet, HttpServletRequest, HttpServletResponse}

import org.eclipse.jetty.server.{Server, ServerConnector}
import org.eclipse.jetty.servlet.{ServletContextHandler, ServletHolder}
import org.eclipse.jetty.util.thread.QueuedThreadPool

class HttpServer(config: SchedulerConfigBase) {
  val threadPool = new QueuedThreadPool(16)
  threadPool.setName("Jetty")

  val server = new Server(threadPool)
  val connector = new ServerConnector(server)
  connector.setPort(config.artifactServerPort)

  val handler = new ServletContextHandler
  handler.addServlet(new ServletHolder(new Servlet()), "/")

  server.setHandler(handler)
  server.addConnector(connector)
  server.start()

  def stop() {
    if (server == null) throw new IllegalStateException("!started")

    server.stop()
    server.join()
  }

  class Servlet extends HttpServlet {
    override def doGet(request: HttpServletRequest, response: HttpServletResponse): Unit = {
      val uri = request.getRequestURI
      if (uri.startsWith("/resource/")) downloadFile(uri.split("/").last, response)
      else if (uri.startsWith("/scale/")) scale(uri.split("/").last.toInt, response)
      else response.sendError(404)
    }

    def downloadFile(file: String, response: HttpServletResponse) {
      response.setContentType("application/zip")
      response.setHeader("Content-Disposition", "attachment; filename=\"" + new File(file).getName + "\"")
      Util.copyAndClose(new FileInputStream(file), response.getOutputStream)
    }

    def scale(scale: Int, response: HttpServletResponse) {
      if (scale >= 0) config.instances = scale
    }
  }

}