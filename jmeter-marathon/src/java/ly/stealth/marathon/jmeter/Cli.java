package ly.stealth.marathon.jmeter;

import joptsimple.OptionException;
import joptsimple.OptionParser;
import joptsimple.OptionSet;
import org.apache.log4j.*;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.util.Arrays;
import java.util.List;

import static java.lang.System.err;
import static java.lang.System.out;

public class Cli {
    public static void main(String[] args) throws Exception {
        initLogging();

        if (args.length == 0) {
            out.println("Usage: start|stop|status");
            System.exit(1);
        }

        String command = args[0];
        args = Arrays.asList(args).subList(1, args.length).toArray(new String[args.length - 1]);

        switch (command) {
            case "status": handleStatus(args); break;
            case "start": handleStart(args); break;
            case "stop": handleStop(args); break;
            default: err.println("Unsupported command " + command); System.exit(1);
        }
    }

    private static void handleStatus(String... args) throws IOException {
        OptionParser parser = new OptionParser();
        parser.accepts("marathon").withRequiredArg().required().ofType(String.class);
        parser.accepts("app").withRequiredArg().ofType(String.class).defaultsTo(Marathon.JmeterServers.DEFAULT_APP);

        OptionSet options = null;
        try {
            options = parser.parse(args);
        } catch (OptionException e) {
            parser.printHelpOn(out);
            err.println(e.getMessage());
            System.exit(1);
        }

        Marathon.url = (String) options.valueOf("marathon");
        String app = (String) options.valueOf("app");
        List<String> endpoints = Marathon.getEndpoints(app);

        if (endpoints.isEmpty()) out.println("App \"" + app + "\" is not started");
        else out.println("App \"" + app + "\"" + " is running\nServers listening on " + Util.join(endpoints, ","));
    }

    private static void handleStart(String... args) throws Exception {
        Marathon.JmeterServers servers = new Marathon.JmeterServers();

        OptionParser parser = new OptionParser();
        parser.accepts("marathon").withRequiredArg().required().ofType(String.class);
        parser.accepts("download").withRequiredArg().required().ofType(String.class);

        parser.accepts("app").withOptionalArg().ofType(String.class).defaultsTo(servers.app);
        parser.accepts("instances").withOptionalArg().ofType(Integer.class).defaultsTo(servers.instances);
        parser.accepts("cpus").withOptionalArg().ofType(Double.class).defaultsTo(servers.cpus);
        parser.accepts("mem").withOptionalArg().ofType(Integer.class).defaultsTo(servers.mem);

        OptionSet options = null;
        try {
            options = parser.parse(args);
        } catch (OptionException e) {
            parser.printHelpOn(out);
            err.println(e.getMessage());
            System.exit(1);
        }

        Marathon.url = (String) options.valueOf("marathon");

        servers.downloadUrl = (String) options.valueOf("download");
        servers.app = (String) options.valueOf("app");
        servers.instances = (int) options.valueOf("instances");

        servers.cpus = (double) options.valueOf("cpus");
        servers.mem = (int) options.valueOf("mem");

        HttpServer server = new HttpServer();
        server.setPort(new URI(servers.downloadUrl).getPort());
        server.setJmeterDistro(jmeterDistro());
        server.start();

        if (Marathon.hasApp(servers.app)) {
            err.println("App \"" + servers.app + "\" is already running");
            System.exit(1);
        }

        out.println("Starting app \"" + servers.app + "\" ...");
        Marathon.startServers(servers);

        out.println("Servers listening on " + Util.join(Marathon.getEndpoints(servers.app), ","));
        server.stop();
    }

    private static void handleStop(String... args) throws IOException {
        OptionParser parser = new OptionParser();
        parser.accepts("marathon").withRequiredArg().required().ofType(String.class);
        parser.accepts("app").withRequiredArg().ofType(String.class).defaultsTo(Marathon.JmeterServers.DEFAULT_APP);

        OptionSet options = null;
        try {
            options = parser.parse(args);
        } catch (OptionException e) {
            parser.printHelpOn(out);
            err.println(e.getMessage());
            System.exit(1);
        }

        Marathon.url = (String) options.valueOf("marathon");
        String app = (String) options.valueOf("app");

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
