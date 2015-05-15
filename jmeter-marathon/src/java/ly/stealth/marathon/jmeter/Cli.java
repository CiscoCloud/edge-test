package ly.stealth.marathon.jmeter;

import org.apache.log4j.*;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static java.lang.System.err;
import static java.lang.System.out;

public class Cli {
    public static void main(String[] args_) throws Exception {
        initLogging();
        List<String> args = new ArrayList<>(Arrays.asList(args_));

        if (args.size() == 0) {
            out.println("Usage: start|stop|status");
            System.exit(1);
        }

        String command = args.get(0);
        args = args.subList(1, args.size());

        switch (command) {
            case "status": handleStatus(args); break;
            case "start": handleStart(args); break;
            case "stop": handleStop(args); break;
            default: err.println("Unsupported command " + command); System.exit(1);
        }
    }

    private static void handleStatus(List<String> args) throws IOException {
        String app = Marathon.JmeterServers.DEFAULT_APP;

        Marathon.url = "http://master:8080";
        List<String> endpoints = Marathon.getEndpoints(app);

        if (endpoints.isEmpty()) out.println("App \"" + app + "\" is not started");
        else out.println("App \"" + app + "\"" + " is running\nServers listening on " + Util.join(endpoints, ","));
    }

    private static void handleStart(List<String> args) throws Exception {
        Marathon.url = "http://master:8080";

        HttpServer server = new HttpServer();
        server.setPort(5000);
        server.setJmeterDistro(jmeterDistro());
        server.start();

        Marathon.JmeterServers servers = new Marathon.JmeterServers();
        servers.downloadUrl = "http://192.168.3.1:5000";
        servers.instances = 2;

        if (Marathon.hasApp(servers.app)) {
            err.println("App \"" + servers.app + "\" is already running");
            System.exit(1);
        }

        out.println("Starting app \"" + servers.app + "\" ...");
        Marathon.startServers(servers);

        out.println("Servers listening on " + Util.join(Marathon.getEndpoints(servers.app), ","));
        server.stop();
    }

    private static void handleStop(List<String> args) throws IOException {
        Marathon.url = "http://master:8080";
        String app = Marathon.JmeterServers.DEFAULT_APP;

        if (!Marathon.hasApp(app)) {
            err.println("App \"" + app + "\" is not started");
            System.exit(1);
        }

        out.println("Stopping app \"" + app + "\" ...");
        Marathon.stopApp(app);
    }

    private static void initLogging() {
        System.setProperty("org.eclipse.jetty.util.log.class", HttpServer.JettyLog4jLogger.class.getName());
        BasicConfigurator.resetConfiguration();

        Logger root = Logger.getRootLogger();
        root.setLevel(Level.INFO);

        Logger.getLogger("org.eclipse.jetty").setLevel(Level.WARN);
        Logger.getLogger("jetty").setLevel(Level.WARN);

        PatternLayout layout = new PatternLayout("%d [%t] %-5p %c %x - %m%n");
        root.addAppender(new ConsoleAppender(layout));
    }

    @SuppressWarnings("ConstantConditions")
    private static File jmeterDistro() {
        String mask = "apache-jmeter.*\\.zip";

        for (File file : new File(".").listFiles())
            if (file.getName().matches(mask)) return file;

        throw new IllegalStateException("No " + mask + " found in .");
    }
}
