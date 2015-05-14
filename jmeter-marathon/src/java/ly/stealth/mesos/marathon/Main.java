package ly.stealth.mesos.marathon;

import org.apache.log4j.*;

import java.io.File;

public class Main {
    private static final Logger logger = Logger.getLogger(Main.class);
    private static File jmeterDistro = jmeterDistro();

    public static void main(String[] args) throws Exception {
        initLogging();

        HttpServer server = new HttpServer();
        server.setPort(5000);
        server.setJmeterDistro(jmeterDistro);
        server.start();

        Agents agents = new Agents();
        agents.setApiUrl("http://192.168.3.1:5000");
        agents.setMarathonUrl("http://master:8080");
        agents.setInstances(3);
        agents.start();

        logger.info("Agents listening on " + Util.join(agents.getEndpoints(), ","));
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
