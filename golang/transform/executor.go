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
	"encoding/json"
	"fmt"
	"github.com/CiscoCloud/edge-test/golang/avro"
	pb "github.com/CiscoCloud/edge-test/golang/proto"
	"github.com/golang/protobuf/proto"
	"github.com/mesos/mesos-go/executor"
	mesos "github.com/mesos/mesos-go/mesosproto"
	kafka "github.com/stealthly/go_kafka_client"
	"io/ioutil"
	"net/http"
	"time"
)

type TransformExecutorConfig struct {
	BrokerList        []string
	SchemaRegistryUrl string
	Topic             string
	Port              int
}

func NewTransformExecutorConfig() *TransformExecutorConfig {
	return new(TransformExecutorConfig)
}

type TransformExecutor struct {
	config   *TransformExecutorConfig
	incoming chan *avro.LogLine
	producer kafka.Producer
	close    chan bool
}

// Creates a new TransformExecutor with a given config.
func NewTransformExecutor(config *TransformExecutorConfig) *TransformExecutor {
	return &TransformExecutor{
		config:   config,
		incoming: make(chan *avro.LogLine),
		close:    make(chan bool),
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
	http.HandleFunc("/", this.handleFunc())

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

func (this *TransformExecutor) handleFunc() func(http.ResponseWriter, *http.Request) {
	avroDecoder := kafka.NewKafkaAvroDecoder(this.config.SchemaRegistryUrl)

	return func(w http.ResponseWriter, r *http.Request) {
		timing := this.timing("received")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		logLine := avro.NewLogLine()
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			{
				err = this.handleJson(body, logLine)
			}
		case "application/x-protobuf":
			{
				err = this.handleProto(body, logLine)
			}
		case "avro/binary":
			{
				err = this.handleAvro(body, logLine, avroDecoder)
			}
		default:
			err = fmt.Errorf("Content-Type %s is invalid", contentType)
		}

		if err != nil {
			fmt.Printf("Got corrupted log line: %s\n", err)
			return
		}

		logLine.Timings = append(logLine.Timings, timing)

		this.incoming <- logLine
	}
}

func (this *TransformExecutor) handleJson(body []byte, logLine *avro.LogLine) error {
	if err := json.Unmarshal(body, logLine); err != nil {
		return err
	}

	// golang's json numbers are always float64's :(
	if logLine.Logtypeid != nil {
		logLine.Logtypeid = int64(logLine.Logtypeid.(float64))
	}

	return nil
}

func (this *TransformExecutor) handleAvro(body []byte, logLine *avro.LogLine, decoder *kafka.KafkaAvroDecoder) error {
	return decoder.DecodeSpecific(body, logLine)
}

func (this *TransformExecutor) handleProto(body []byte, logLine *avro.LogLine) error {
	protoLogLine := &pb.LogLine{}
	if err := proto.Unmarshal(body, protoLogLine); err != nil {
		return err
	}

	this.protoToAvroLogLine(protoLogLine, logLine)
	return nil
}

func (this *TransformExecutor) protoToAvroLogLine(protoLogLine *pb.LogLine, logLine *avro.LogLine) *avro.LogLine {
	logLine.Line = protoLogLine.GetLine()
	logLine.Source = protoLogLine.GetSource()
	logLine.Logtypeid = protoLogLine.GetLogtypeid()

	logLine.Tag = make(map[string]string)
	for _, tag := range protoLogLine.GetTag() {
		logLine.Tag[tag.GetKey()] = tag.GetValue()
	}

	logLine.Timings = make([]*avro.KV, 0)
	for _, timing := range protoLogLine.GetTimings() {
		kv := avro.NewKV()
		kv.Key = timing.GetKey()
		kv.Value = timing.GetValue()
		logLine.Timings = append(logLine.Timings, kv)
	}

	return logLine
}

func (this *TransformExecutor) timing(name string) *avro.KV {
	timing := avro.NewKV()
	timing.Key = name
	timing.Value = time.Now().Unix() * 1000
	return timing
}
