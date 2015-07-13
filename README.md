# edge-test
Test implementation for containerization of the edge components and their validation results.

Components:

* [Common](https://github.com/CiscoCloud/edge-test/tree/master/common/) - contains all base Mesos scheduler/executor logic, avro/json/proto -> avro transformation, LogLine model.
* Frameworks: 

[Dropwizard](https://github.com/CiscoCloud/edge-test/tree/master/dropwizard)   
[Finagle](https://github.com/CiscoCloud/edge-test/tree/master/finagle)   
[Golang](https://github.com/CiscoCloud/edge-test/tree/master/golang)   
[Netty](https://github.com/CiscoCloud/edge-test/tree/master/netty)   
[Play](https://github.com/CiscoCloud/edge-test/tree/master/play)   
[Spray](https://github.com/CiscoCloud/edge-test/tree/master/spray)   
[Unfiltered](https://github.com/CiscoCloud/edge-test/tree/master/unfiltered)   

* [Go-perf](https://github.com/CiscoCloud/edge-test/tree/master/go-perf) - simple go tool to create load simulating multiple concurrent connections.
* [Scripts](https://github.com/CiscoCloud/edge-test/tree/master/scripts) - simple bash scripts to create load (using curl).
* [Vagrant](https://github.com/CiscoCloud/edge-test/tree/master/vagrant) - contains vagrant scripts to spin up Mesos locally.