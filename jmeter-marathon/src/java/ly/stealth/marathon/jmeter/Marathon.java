package ly.stealth.marathon.jmeter;

import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.JSONValue;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

public class Marathon {
    public static String url = "http://master:8080";

    public static boolean hasApp(String id) throws IOException {
        return getApp(id) != null;
    }

    public static List<String> getEndpoints(String id) throws IOException {
        @SuppressWarnings("unchecked") List<JSONObject> tasks = getTasks(id);
        if (tasks == null) return Collections.emptyList();

        List<String> endpoints = new ArrayList<>();
        for (JSONObject task : tasks)
            endpoints.add(task.get("host") + ":" + ((JSONArray)task.get("ports")).get(0));

        return endpoints;
    }

    public static void startServers(JmeterServers servers) throws IOException, InterruptedException {
        sendRequest("/v2/apps", "POST", servers.appJson());

        for(;;) {
            @SuppressWarnings("unchecked")
            List<JSONObject> tasks = getTasks(servers.id);
            if (tasks == null) throw new IllegalStateException("App " + servers.id + " not found");

            int started = 0;
            for (JSONObject task : tasks)
                if (task.get("startedAt") != null) started ++;

            if (started == servers.instances) break;
            Thread.sleep(1000);
        }
    }

    public static void stopServers(String id) throws IOException {
        sendRequest("/v2/apps/" + id, "DELETE", null);
    }

    private static JSONArray getTasks(String id) throws IOException {
        JSONObject app = getApp(id);
        return app != null ? (JSONArray) app.get("tasks") : null;
    }

    private static JSONObject getApp(String id) throws IOException {
        JSONObject response = sendRequest("/v2/apps/" + id, "GET", null);
        if (response == null) return null;
        return (JSONObject) response.get("app");
    }

    private static JSONObject sendRequest(String uri, String method, JSONObject json) throws IOException {
        URL url = new URL(Marathon.url + uri);
        HttpURLConnection c = (HttpURLConnection) url.openConnection();
        try {
            c.setRequestMethod(method);

            if (method.equalsIgnoreCase("POST")) {
                byte[] body = json.toString().getBytes("utf-8");
                c.setDoOutput(true);
                c.setRequestProperty("Content-Type", "application/json");
                c.setRequestProperty("Content-Length", "" + body.length);
                c.getOutputStream().write(body);
            }

            return (JSONObject) JSONValue.parse(new InputStreamReader(c.getInputStream(), "utf-8"));
        } catch (IOException e) {
            if (c.getResponseCode() == 404 && method.equals("GET"))
                return null;

            ByteArrayOutputStream response = new ByteArrayOutputStream();
            InputStream err = c.getErrorStream();
            if (err == null) throw e;

            Util.copyAndClose(err, response);
            IOException ne = new IOException(e.getMessage() + ": " + response.toString("utf-8"));
            ne.setStackTrace(e.getStackTrace());
            throw ne;
        } finally {
            c.disconnect();
        }
    }

    public static class JmeterServers {
        public static final String DEFAULT_ID = "jmeter";

        public String downloadUrl;
        public String id = DEFAULT_ID;
        public int instances = 1;

        @SuppressWarnings("unchecked")
        private JSONObject appJson() {
            JSONObject obj = new JSONObject();

            obj.put("id", id);
            obj.put("cpus", 1);
            obj.put("mem", 128);

            obj.put("instances", instances);
            obj.put("ports", Arrays.asList(0));
            obj.put("uris", Arrays.asList(downloadUrl + "/jmeter/apache-jmeter.zip", downloadUrl + "/jmeter/agent.sh"));

            String startJmeter = "./jmeter.sh -s -Jserver.rmi.localport=$PORT0 -Jserver.rmi.port=$PORT";
            obj.put("cmd", "cd apache-jmeter*/bin && chmod +x *.sh && " + startJmeter);

            return obj;
        }
    }
}
