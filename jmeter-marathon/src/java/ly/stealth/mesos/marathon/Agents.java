package ly.stealth.mesos.marathon;

import org.json.simple.JSONObject;
import org.json.simple.JSONValue;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Arrays;

public class Agents {
    public static final File script;
    static {
        try { script = tmpScript(); }
        catch (IOException e) { throw new Error(e); }
    }

    private String apiUrl = "http://localhost:5000";
    private String marathonUrl = "http://master:8080";
    private int instances = 1;


    public String getApiUrl() { return apiUrl; }
    public void setApiUrl(String apiUrl) { this.apiUrl = apiUrl; }

    public String getMarathonUrl() { return marathonUrl; }
    public void setMarathonUrl(String marathonUrl) { this.marathonUrl = marathonUrl; }

    public int getInstances() { return instances; }
    public void setInstances(int instances) { this.instances = instances; }


    public void submit() throws IOException {
        JSONObject response = sendRequest("/v2/apps", "POST", appJson());
        System.out.println(response);
    }

    @SuppressWarnings("unchecked")
    private JSONObject appJson() {
        JSONObject obj = new JSONObject();

        obj.put("id", "jmeter-agent");
        obj.put("cpus", 1);
        obj.put("mem", 128);

        obj.put("instances", instances);
        obj.put("ports", Arrays.asList(0));

        obj.put("uris", Arrays.asList(apiUrl + "/jmeter/apache-jmeter.zip", apiUrl + "/jmeter/agent.sh"));
        obj.put("cmd", "chmod +x agent.sh && ./agent.sh");

        return obj;
    }

    private JSONObject sendRequest(String uri, String method, JSONObject json) throws IOException {
        URL url = new URL(marathonUrl + uri);
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

    private static File tmpScript() throws IOException {
        InputStream stream = Agents.class.getResourceAsStream("agent.sh");
        if (stream == null) throw new IllegalStateException("No agent.sh on classpath near " + Agents.class.getName());

        File file = File.createTempFile("agent.sh", null);
        file.deleteOnExit();

        Util.copyAndClose(stream, new FileOutputStream(file));
        return file;
    }
}
