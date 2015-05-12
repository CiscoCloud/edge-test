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
    "github.com/mesos/mesos-go/executor"
    "fmt"
    "flag"
    "github.com/CiscoCloud/edge-test/golang/transform"
)

func parseAndValidateExecutorArgs() {
    flag.Parse()
}

func main() {
    parseAndValidateExecutorArgs()
    fmt.Println("Starting Transform Executor")

    driverConfig := executor.DriverConfig {
        Executor: transform.NewTransformExecutor(),
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
