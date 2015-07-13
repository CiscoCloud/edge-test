# edge-test
Test implementation for containerization of the edge components and their validation results.

Components:

* Common - contains all base Mesos scheduler/executor logic, avro/json/proto -> avro transformation, LogLine model.
* Frameworks: 

Dropwizard   
Finagle   
Golang   
Netty   
Play   
Spray   
Unfiltered   

* Go-perf - simple go tool to create load simulating multiple concurrent connections.
* Scripts - simple bash scripts to create load (using curl).
* Vagrant - contains vagrant scripts to spin up Mesos locally.