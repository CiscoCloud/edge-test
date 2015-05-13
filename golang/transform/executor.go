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

package transform

import (
    "github.com/mesos/mesos-go/executor"
    mesos "github.com/mesos/mesos-go/mesosproto"
    "fmt"
    "net/http"
    kafka "github.com/stealthly/go_kafka_client"
    "github.com/CiscoCloud/edge-test/golang/avro"
    "io/ioutil"
    "encoding/json"
    "time"
)

type TransformExecutorConfig struct {
    Format string
    BrokerList []string
    SchemaRegistryUrl string
    Topic string
    Port int
}

func NewTransformExecutorConfig() *TransformExecutorConfig {
    return &TransformExecutorConfig {
        Format: "json",
    }
}

type TransformExecutor struct {
    config *TransformExecutorConfig
    incoming chan *avro.LogLine
    producer kafka.Producer
    close chan bool
}

// Creates a new TransformExecutor with a given config.
func NewTransformExecutor(config *TransformExecutorConfig) *TransformExecutor {
    return &TransformExecutor{
        config: config,
        incoming: make(chan *avro.LogLine),
        close: make(chan bool),
    }
}

// mesos.Executor interface method.
// Invoked once the executor driver has been able to successfully connect with Mesos.
// Not used by TransformExecutor yet.
func (this *TransformExecutor) Registered(driver executor.ExecutorDriver, execInfo *mesos.ExecutorInfo, fwinfo *mesos.FrameworkInfo, slaveInfo *mesos.SlaveInfo) {
    fmt.Printf("Registered Executor on slave %s\n", slaveInfo.GetHostname())
}

// mesos.Executor interface method.
// Invoked when the executor re-registers with a restarted slave.
func (this *TransformExecutor) Reregistered(driver executor.ExecutorDriver, slaveInfo *mesos.SlaveInfo) {
    fmt.Printf("Re-registered Executor on slave %s\n", slaveInfo.GetHostname())
}

// mesos.Executor interface method.
// Invoked when the executor becomes "disconnected" from the slave.
func (this *TransformExecutor) Disconnected(executor.ExecutorDriver) {
    fmt.Println("Executor disconnected.")
}

// mesos.Executor interface method.
// Invoked when a task has been launched on this executor.
func (this *TransformExecutor) LaunchTask(driver executor.ExecutorDriver, taskInfo *mesos.TaskInfo) {
    fmt.Printf("Launching task %s with command %s\n", taskInfo.GetName(), taskInfo.Command.GetValue())

    runStatus := &mesos.TaskStatus{
        TaskId: taskInfo.GetTaskId(),
        State:  mesos.TaskState_TASK_RUNNING.Enum(),
    }

    if _, err := driver.SendStatusUpdate(runStatus); err != nil {
        fmt.Printf("Failed to send status update: %s\n", runStatus)
    }

    go func() {
        this.startHTTPServer()
        this.startProducer()
        <-this.close
        close(this.incoming)

        // finish task
        fmt.Printf("Finishing task %s\n", taskInfo.GetName())
        finStatus := &mesos.TaskStatus{
            TaskId: taskInfo.GetTaskId(),
            State:  mesos.TaskState_TASK_FINISHED.Enum(),
        }
        if _, err := driver.SendStatusUpdate(finStatus); err != nil {
            fmt.Printf("Failed to send status update: %s\n", finStatus)
        }
        fmt.Printf("Task %s has finished\n", taskInfo.GetName())
    }()
}

// mesos.Executor interface method.
// Invoked when a task running within this executor has been killed.
func (this *TransformExecutor) KillTask(_ executor.ExecutorDriver, taskId *mesos.TaskID) {
    fmt.Println("Kill task")

    select {
    case this.close <- true:
    default:
    }
}

// mesos.Executor interface method.
// Invoked when a framework message has arrived for this executor.
func (this *TransformExecutor) FrameworkMessage(driver executor.ExecutorDriver, msg string) {
    fmt.Printf("Got framework message: %s\n", msg)
}

// mesos.Executor interface method.
// Invoked when the executor should terminate all of its currently running tasks.
func (this *TransformExecutor) Shutdown(executor.ExecutorDriver) {
    fmt.Println("Shutting down the executor")

    select {
    case this.close <- true:
    default:
    }
}

// mesos.Executor interface method.
// Invoked when a fatal error has occured with the executor and/or executor driver.
func (this *TransformExecutor) Error(driver executor.ExecutorDriver, err string) {
    fmt.Printf("Got error message: %s\n", err)
}

func (this *TransformExecutor) startHTTPServer() {
    var handleFunc func (http.ResponseWriter, *http.Request)
    switch this.config.Format {
        case "json": handleFunc = this.jsonHandleFunc
        case "avro": handleFunc = this.avroHandleFunc
        case "proto": handleFunc = this.protoHandleFunc
    }

    http.HandleFunc("/", handleFunc)

    go http.ListenAndServe(fmt.Sprintf(":%d", this.config.Port), nil)
}

func (this *TransformExecutor) startProducer() {
    producerConfig := kafka.DefaultProducerConfig()

    producerConfig.KeyEncoder = kafka.NewKafkaAvroEncoder(this.config.SchemaRegistryUrl)
    producerConfig.ValueEncoder = producerConfig.KeyEncoder
    producerConfig.BrokerList = this.config.BrokerList

    this.producer = kafka.NewSaramaProducer(producerConfig)
    go this.produceRoutine()
}

func (this *TransformExecutor) produceRoutine() {
    for msg := range this.incoming {
        msg.Timings = append(msg.Timings, this.timing("sent"))
        this.producer.Input() <- &kafka.ProducerMessage{Topic: this.config.Topic, Value: msg}
    }
}

func (this *TransformExecutor) jsonHandleFunc(w http.ResponseWriter, r *http.Request) {
    timing := this.timing("received")
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    logLine := avro.NewLogLine()
    if err = json.Unmarshal(body, logLine); err != nil {
        panic(err)
    }

    // golang's json numbers are always float64's :(
    if logLine.Logtypeid != nil {
        fmt.Println("changing logtypeid type")
        logLine.Logtypeid = int64(logLine.Logtypeid.(float64))
    }

    logLine.Timings = append(logLine.Timings, timing)

    this.incoming <- logLine
}

func (this *TransformExecutor) avroHandleFunc(w http.ResponseWriter, r *http.Request) {
    //TODO
}

func (this *TransformExecutor) protoHandleFunc(w http.ResponseWriter, r *http.Request) {
    //TODO
}

func (this *TransformExecutor) timing(name string) *avro.KV {
    timing := avro.NewKV()
    timing.Key = "received"
    timing.Value = time.Now().Unix() * 1000
    return timing
}
