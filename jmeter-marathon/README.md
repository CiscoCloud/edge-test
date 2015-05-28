# Intro
JMeter Marathon is a tool intended to run JMeter on Mesos/Marathon cluster.

# Quick start
Running `./gradlew jar` will produce `jmeter-marathon*.jar` into project folder.
After that you can use `./jmeter-marathon.sh` script to run commands via CLI.

Example:
```
# ./jmeter-marathon.sh start --marathon=http://master:8080 --download=http://192.168.3.1:5001 --instances=2
Starting app "jmeter" ...
Servers listening on slave0:31541,master:31053

```

Now we can use JMeter GUI or console client with our running JMeter servers by passing `-J"remote_hosts=slave0:31541,master:31053"`.
After that the test plan could be executed on the specified remote servers.

# CLI usage
This project provides CLI with following commands:
```
# ./jmeter-marathon.sh help
Usage: help {cmd}|start|stop|status

```

Start/stop starts and stops JMeter servers on Mesos cluster starting/destroying Marathon app.
Status provides a status for Marathon app.

For more information about particular command, please refer to correspondent help. Example:
```
# ./jmeter-marathon.sh help stop
Stop servers
Usage: stop [options]

Option (* = required)  Description
---------------------  -----------
--app                  marathon app id (default: jmeter)
* --marathon           marathon url (http://master:8080)

```

