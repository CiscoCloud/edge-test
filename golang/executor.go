// +build executor

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
	"flag"
	"fmt"
	"github.com/CiscoCloud/edge-test/golang/transform"
	"github.com/mesos/mesos-go/executor"
	"os"
	"strings"
)

var port = flag.Int("port", 0, "Port to bind to")
var schemaRegistryUrl = flag.String("schema.registry", "", "Avro Schema Registry url.")
var brokerList = flag.String("broker.list", "", "Comma separated list of brokers for producer.")
var topic = flag.String("topic", "", "Topic to produce transformed data to.")

func parseAndValidateExecutorArgs() {
	flag.Parse()

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

	if *port == 0 {
		fmt.Println("port flag is required.")
		os.Exit(1)
	}
}

func main() {
	parseAndValidateExecutorArgs()
	fmt.Println("Starting Transform Executor")

	executorConfig := transform.NewTransformExecutorConfig()
	executorConfig.SchemaRegistryUrl = *schemaRegistryUrl
	executorConfig.BrokerList = strings.Split(*brokerList, ",")
	executorConfig.Topic = *topic
	executorConfig.Port = *port

	driverConfig := executor.DriverConfig{
		Executor: transform.NewTransformExecutor(executorConfig),
	}
	driver, err := executor.NewMesosExecutorDriver(driverConfig)

	if err != nil {
		fmt.Println("Unable to create a ExecutorDriver ", err.Error())
	}

	_, err = driver.Start()
	if err != nil {
		fmt.Println("Got error:", err)
		return
	}
	fmt.Println("Executor process has started and running.")
	driver.Join()
}
