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

import io.netty.bootstrap.ServerBootstrap
import io.netty.buffer.Unpooled
import io.netty.channel.nio.NioEventLoopGroup
import io.netty.channel.socket.SocketChannel
import io.netty.channel.socket.nio.NioServerSocketChannel
import io.netty.channel.{ChannelFutureListener, ChannelHandlerContext, ChannelInitializer, SimpleChannelInboundHandler}
import io.netty.handler.codec.http.HttpHeaders.Names._
import io.netty.handler.codec.http.HttpResponseStatus._
import io.netty.handler.codec.http.HttpVersion._
import io.netty.handler.codec.http._
import io.netty.handler.logging.{LogLevel, LoggingHandler}
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
  private val bossGroup = new NioEventLoopGroup(1)
  private val workerGroup = new NioEventLoopGroup()

  private val transformer = new Transform(config)

  val serverBootstrap = new ServerBootstrap()
  serverBootstrap.group(bossGroup, workerGroup)
    .channel(classOf[NioServerSocketChannel])
    .handler(new LoggingHandler(LogLevel.INFO))
    .childHandler(new ServerInitializer(config, transformer))

  new Thread() {
    override def run() {
      try {
        serverBootstrap.bind(config.port).sync().channel().closeFuture().sync()
      } catch {
        case e: Throwable => e.printStackTrace()
      }
    }
  }.start()
}

class ServerInitializer(config: ExecutorConfigBase, transformer: Transform) extends ChannelInitializer[SocketChannel] {
  override def initChannel(ch: SocketChannel) {
    val pipeline = ch.pipeline()
    pipeline.addLast(new HttpRequestDecoder)
    pipeline.addLast(new HttpResponseEncoder)
    pipeline.addLast(new ServerHandler(config, transformer))
  }
}

class ServerHandler(config: ExecutorConfigBase, transformer: Transform) extends SimpleChannelInboundHandler[Any] {
  private var request: HttpRequest = null
  private var contentType: String = ""

  override def channelReadComplete(ctx: ChannelHandlerContext) {
    ctx.flush()
  }

  override protected def channelRead0(ctx: ChannelHandlerContext, msg: Any) {
    msg match {
      case request: HttpRequest =>
        this.request = request

        if (HttpHeaders.is100ContinueExpected(request)) {
          send100Continue(ctx)
        }
        contentType = request.headers().get("Content-Type")
      case _ =>
    }

    msg match {
      case content: HttpContent =>
        if (content.content().isReadable) {
          val body = new Array[Byte](content.content().capacity())
          content.content().getBytes(0, body)
          if (!config.sync) {
            new Thread {
              override def run() {
                transformer.transform(body, contentType, "Netty")
              }
            }.start()
          } else transformer.transform(body, contentType, "Netty")
        }

        if (msg.isInstanceOf[LastHttpContent]) {
          if (!writeResponse(content, ctx)) {
            ctx.writeAndFlush(Unpooled.EMPTY_BUFFER).addListener(ChannelFutureListener.CLOSE)
          }
        }
      case _ =>
    }
  }

  private def send100Continue(ctx: ChannelHandlerContext) {
    val response = new DefaultFullHttpResponse(HTTP_1_1, CONTINUE)
    ctx.write(response)
  }

  private def writeResponse(current: HttpObject, ctx: ChannelHandlerContext): Boolean = {
    val keepAlive = HttpHeaders.isKeepAlive(request)

    val status = if (current.getDecoderResult.isSuccess) OK else BAD_REQUEST
    val response = new DefaultFullHttpResponse(HTTP_1_1, status)

    response.headers().set(CONTENT_TYPE, "text/plain; charset=UTF-8")

    if (keepAlive) {
      response.headers().set(CONTENT_LENGTH, "0")
      response.headers().set(CONNECTION, HttpHeaders.Values.KEEP_ALIVE)
    }

    ctx.write(response)

    keepAlive
  }

  override def exceptionCaught(ctx: ChannelHandlerContext, cause: Throwable) {
    cause.printStackTrace()
    ctx.close()
  }
}

