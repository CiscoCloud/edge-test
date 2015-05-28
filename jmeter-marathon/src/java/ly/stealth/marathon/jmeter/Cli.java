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

        try { handle(args); }
        catch (Error error) {
            out.println();
            err.println("Error: " + Util.uncapitalize(error.getMessage()));
            System.exit(1);
        }
    }

    private static void handle(String... args) throws Exception {
        if (args.length == 0) {
            handleHelp(null);
            throw new Error("Command required");
        }

        String command = args[0];
        args = Arrays.asList(args).subList(1, args.length).toArray(new String[args.length - 1]);

        switch (command) {
            case "help": handleHelp(args.length > 0 ? args[0] : null); break;
            case "status": handleStatus(args, false); break;
            case "start": handleStart(args, false); break;
            case "stop": handleStop(args, false); break;
            default: handleHelp(null); throw new Error("Unsupported command");
        }
    }

    private static void handleHelp(String command) throws Exception {
        if (command == null) {
            out.println("Usage: help {cmd}|start|stop|status");
            return;
        }

        switch (command) {
            case "help": out.println("Print general or command-specific help\nUsage: help {command}"); break;
            case "status": handleStatus(null, true); break;
            case "start": handleStart(null, true); break;
            case "stop": handleStop(null, true); break;
            default: handleHelp(null); throw new Error("Unsupported command");
        }

    }

    private static void handleStatus(String[] args, boolean help) throws IOException {
        OptionParser parser = new OptionParser();
        parser.accepts("marathon", "marathon url (http://master:8080)").withRequiredArg().required().ofType(String.class);
        parser.accepts("app", "marathon app id").withRequiredArg().ofType(String.class).defaultsTo(Marathon.JmeterServers.DEFAULT_APP);

        if (help) {
            out.println("Print servers status\nUsage: status [options]\n");
            parser.printHelpOn(out);
            return;
        }

        OptionSet options;
        try {
            options = parser.parse(args);
        } catch (OptionException e) {
            parser.printHelpOn(out);
            throw new Error(e.getMessage());
        }

        Marathon.url = (String) options.valueOf("marathon");
        String app = (String) options.valueOf("app");
        List<String> endpoints = Marathon.getEndpoints(app);

        if (endpoints.isEmpty()) out.println("App \"" + app + "\" is not started");
        else out.println("App \"" + app + "\"" + " is running\nServers listening on " + Util.join(endpoints, ","));
    }

    private static void handleStart(String[] args, boolean help) throws Exception {
        Marathon.JmeterServers servers = new Marathon.JmeterServers();

        OptionParser parser = new OptionParser();
        parser.accepts("marathon", "marathon url (http://master:8080)").withRequiredArg().required().ofType(String.class);
        parser.accepts("download", "url to download jmeter (http://192.168.3.1:5000)").withRequiredArg().required().ofType(String.class);

        parser.accepts("app", "marathon app id").withOptionalArg().ofType(String.class).defaultsTo(servers.app);
        parser.accepts("instances", "number of servers").withOptionalArg().ofType(Integer.class).defaultsTo(servers.instances);
        parser.accepts("cpus", "amount of cpu to use").withOptionalArg().ofType(Double.class).defaultsTo(servers.cpus);
        parser.accepts("mem", "amount of memory to use").withOptionalArg().ofType(Integer.class).defaultsTo(servers.mem);

        if (help) {
            out.println("Start servers\nUsage: start [options]\n");
            parser.printHelpOn(out);
            return;
        }

        OptionSet options;
        try {
            options = parser.parse(args);
        } catch (OptionException e) {
            parser.printHelpOn(out);
            throw new Error(e.getMessage());
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

        if (Marathon.hasApp(servers.app))
            throw new Error("App \"" + servers.app + "\" is already running");

        out.println("Starting app \"" + servers.app + "\" ...");
        Marathon.startServers(servers);

        out.println("Servers listening on " + Util.join(Marathon.getEndpoints(servers.app), ","));
        server.stop();
    }

    private static void handleStop(String[] args, boolean help) throws IOException {
        OptionParser parser = new OptionParser();
        parser.accepts("marathon", "marathon url (http://master:8080)").withRequiredArg().required().ofType(String.class);
        parser.accepts("app", "marathon app id").withRequiredArg().ofType(String.class).defaultsTo(Marathon.JmeterServers.DEFAULT_APP);

        if (help) {
            out.println("Stop servers\nUsage: stop [options]\n");
            parser.printHelpOn(out);
            return;
        }

        OptionSet options;
        try {
            options = parser.parse(args);
        } catch (OptionException e) {
            parser.printHelpOn(out);
            throw new Error(e.getMessage());
        }

        Marathon.url = (String) options.valueOf("marathon");
        String app = (String) options.valueOf("app");

        if (!Marathon.hasApp(app))
            throw new Error("App \"" + app + "\" is not started");

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

        throw new IllegalStateException("No " + mask + " found in . folder");
    }

    public static class Error extends java.lang.Error {
        public Error(String message) { super(message); }
    }
}
