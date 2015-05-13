// +build framework

/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package main

import (
    "github.com/mesos/mesos-go/mesosproto"
    "github.com/mesos/mesos-go/scheduler"
    "github.com/golang/protobuf/proto"
    "net/http"
    "flag"
    "fmt"
    "os"
    "strings"
    "os/signal"
    "github.com/CiscoCloud/edge-test/golang/transform"
)

var instances = flag.Int("instances", 1, "Number of tasks to run.")
var artifactServerHost = flag.String("artifact.host", "master", "Binding host for artifact server.")
var artifactServerPort = flag.Int("artifact.port", 6666, "Binding port for artifact server.")
var master = flag.String("master", "127.0.0.1:5050", "Mesos Master address <ip:port>.")
var executorArchiveName = flag.String("executor.archive", "executor.zip", "Executor archive name. Absolute or relative path are both ok.")
var executorBinaryName = flag.String("executor.name", "executor", "Executor binary name contained in archive.")
var format = flag.String("format", "json", "Format of messages to expect. 'json', 'avro' or 'proto'")
var schemaRegistryUrl = flag.String("schema.registry", "", "Avro Schema Registry url.")
var brokerList = flag.String("broker.list", "", "Comma separated list of brokers for producer.")
var topic = flag.String("topic", "", "Topic to produce transformed data to.")

func parseAndValidateSchedulerArgs() {
    flag.Parse()

    switch *format {
        case "json", "avro", "proto":
        default: {
            fmt.Println("Invalid format specified.")
            os.Exit(1)
        }
    }

    if *schemaRegistryUrl == "" {
        fmt.Println("schema.registry flag is required.")
        os.Exit(1)
    }

    if *brokerList == "" {
        fmt.Println("broker.list flag is required.")
        os.Exit(1)
    }

    if *topic == "" {
        fmt.Println("topic flag is required.")
        os.Exit(1)
    }
}

func startArtifactServer() {
    //if the full path is given, take the last token only
    path := strings.Split(*executorArchiveName, "/")
    http.HandleFunc(fmt.Sprintf("/%s", path[len(path)-1]), func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, *executorArchiveName)
    })
    http.ListenAndServe(fmt.Sprintf("%s:%d", *artifactServerHost, *artifactServerPort), nil)
}

func main() {
    parseAndValidateSchedulerArgs()

    ctrlc := make(chan os.Signal, 1)
    signal.Notify(ctrlc, os.Interrupt)

    go startArtifactServer()

    frameworkInfo := &mesosproto.FrameworkInfo{
        User: proto.String(""),
        Name: proto.String("Go LogLine Transform Framework"),
    }

    schedulerConfig := transform.NewTransformSchedulerConfig()
    schedulerConfig.ArtifactServerHost = *artifactServerHost
    schedulerConfig.ArtifactServerPort = *artifactServerPort
    schedulerConfig.ExecutorArchiveName = *executorArchiveName
    schedulerConfig.ExecutorBinaryName = *executorBinaryName
    schedulerConfig.Instances = *instances
    schedulerConfig.Format = *format
    schedulerConfig.SchemaRegistryUrl = *schemaRegistryUrl
    schedulerConfig.BrokerList = *brokerList
    schedulerConfig.Topic = *topic

    transformScheduler := transform.NewTransformScheduler(schedulerConfig)
    driverConfig := scheduler.DriverConfig{
        Scheduler: transformScheduler,
        Framework: frameworkInfo,
        Master: *master,
    }

    driver, err := scheduler.NewMesosSchedulerDriver(driverConfig)
    go func() {
        <-ctrlc
        transformScheduler.Shutdown(driver)
        driver.Stop(false)
    }()

    if err != nil {
        fmt.Println("Unable to create a SchedulerDriver ", err.Error())
    }

    if stat, err := driver.Run(); err != nil {
        fmt.Println("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
    }
}
