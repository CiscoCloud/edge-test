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
	"fmt"
	"github.com/golang/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	"github.com/mesos/mesos-go/scheduler"
	"strings"
	"sync/atomic"
)

type TransformSchedulerConfig struct {
	// Number of CPUs allocated for each created Mesos task.
	CpuPerTask float64

	// Number of RAM allocated for each created Mesos task.
	MemPerTask float64

	// Artifact server host name. Will be used to fetch the executor.
	ArtifactServerHost string

	// Artifact server port.Will be used to fetch the executor.
	ArtifactServerPort int

	// Name of the executor archive file.
	ExecutorArchiveName string

	// Name of the executor binary file contained in the executor archive.
	ExecutorBinaryName string

	// Maximum retries to kill a task.
	KillTaskRetries int

	// Number of task instances to run.
	Instances int

	// Producer config file name.
	ProducerConfig string

	// Topic to produce transformed data to.
	Topic string

	// Flag to respond only after decoding-encoding is done.
	Sync bool
}

func NewTransformSchedulerConfig() *TransformSchedulerConfig {
	return &TransformSchedulerConfig{
		CpuPerTask:      0.2,
		MemPerTask:      256,
		KillTaskRetries: 3,
	}
}

type TransformScheduler struct {
	config           *TransformSchedulerConfig
	runningInstances int32
	tasks            []*mesos.TaskID
}

func NewTransformScheduler(config *TransformSchedulerConfig) *TransformScheduler {
	return &TransformScheduler{
		config: config,
	}
}

// mesos.Scheduler interface method.
// Invoked when the scheduler successfully registers with a Mesos master.
func (this *TransformScheduler) Registered(driver scheduler.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	fmt.Printf("Framework Registered with Master %s\n", masterInfo)
}

// mesos.Scheduler interface method.
// Invoked when the scheduler re-registers with a newly elected Mesos master.
func (this *TransformScheduler) Reregistered(driver scheduler.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	fmt.Printf("Framework Re-Registered with Master %s\n", masterInfo)
}

// mesos.Scheduler interface method.
// Invoked when the scheduler becomes "disconnected" from the master.
func (this *TransformScheduler) Disconnected(scheduler.SchedulerDriver) {
	fmt.Println("Disconnected")
}

// mesos.Scheduler interface method.
// Invoked when resources have been offered to this framework.
func (this *TransformScheduler) ResourceOffers(driver scheduler.SchedulerDriver, offers []*mesos.Offer) {
	fmt.Println("Received offers")

	if int(this.runningInstances) > this.config.Instances {
		toKill := int(this.runningInstances) - this.config.Instances
		for i := 0; i < toKill; i++ {
			driver.KillTask(this.tasks[i])
		}

		this.tasks = this.tasks[toKill:]
	}

	offersAndTasks := make(map[*mesos.Offer][]*mesos.TaskInfo)
	for _, offer := range offers {
		cpus := getScalarResources(offer, "cpus")
		mems := getScalarResources(offer, "mem")
		ports := getRangeResources(offer, "ports")

		remainingCpus := cpus
		remainingMems := mems

		var tasks []*mesos.TaskInfo
		for int(this.getRunningInstances()) < this.config.Instances && this.config.CpuPerTask <= remainingCpus && this.config.MemPerTask <= remainingMems && len(ports) > 0 {
			port := this.takePort(&ports)
			taskPort := &mesos.Value_Range{Begin: port, End: port}
			taskId := &mesos.TaskID{
				Value: proto.String(fmt.Sprintf("golang-%s-%d", *offer.Hostname, *port)),
			}

			task := &mesos.TaskInfo{
				Name:     proto.String(taskId.GetValue()),
				TaskId:   taskId,
				SlaveId:  offer.SlaveId,
				Executor: this.createExecutor(this.getRunningInstances(), *port),
				Resources: []*mesos.Resource{
					util.NewScalarResource("cpus", float64(this.config.CpuPerTask)),
					util.NewScalarResource("mem", float64(this.config.MemPerTask)),
					util.NewRangesResource("ports", []*mesos.Value_Range{taskPort}),
				},
			}
			fmt.Printf("Prepared task: %s with offer %s for launch. Ports: %s\n", task.GetName(), offer.Id.GetValue(), taskPort)

			tasks = append(tasks, task)
			remainingCpus -= this.config.CpuPerTask
			remainingMems -= this.config.MemPerTask
			ports = ports[1:]

			this.tasks = append(this.tasks, taskId)
			this.incRunningInstances()
		}
		fmt.Printf("Launching %d tasks for offer %s\n", len(tasks), offer.Id.GetValue())
		offersAndTasks[offer] = tasks
	}

	unlaunchedTasks := this.config.Instances - int(this.getRunningInstances())
	if unlaunchedTasks > 0 {
		fmt.Printf("There are still %d tasks to be launched and no more resources are available.", unlaunchedTasks)
	}

	for _, offer := range offers {
		tasks := offersAndTasks[offer]
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
	}
}

// mesos.Scheduler interface method.
// Invoked when the status of a task has changed.
func (this *TransformScheduler) StatusUpdate(driver scheduler.SchedulerDriver, status *mesos.TaskStatus) {
	fmt.Printf("Status update: task %s is in state %s\n", status.TaskId.GetValue(), status.State.Enum().String())

	if status.GetState() == mesos.TaskState_TASK_LOST || status.GetState() == mesos.TaskState_TASK_FAILED || status.GetState() == mesos.TaskState_TASK_FINISHED {
		this.removeTask(status.GetTaskId())
		this.decRunningInstances()
	}
}

// mesos.Scheduler interface method.
// Invoked when an offer is no longer valid.
func (this *TransformScheduler) OfferRescinded(scheduler.SchedulerDriver, *mesos.OfferID) {}

// mesos.Scheduler interface method.
// Invoked when an executor sends a message.
func (this *TransformScheduler) FrameworkMessage(scheduler.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, string) {
}

// mesos.Scheduler interface method.
// Invoked when a slave has been determined unreachable
func (this *TransformScheduler) SlaveLost(scheduler.SchedulerDriver, *mesos.SlaveID) {}

// mesos.Scheduler interface method.
// Invoked when an executor has exited/terminated.
func (this *TransformScheduler) ExecutorLost(scheduler.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, int) {
}

// mesos.Scheduler interface method.
// Invoked when there is an unrecoverable error in the scheduler or scheduler driver.
func (this *TransformScheduler) Error(driver scheduler.SchedulerDriver, err string) {
	fmt.Printf("Scheduler received error: %s\n", err)
}

// Gracefully shuts down all running tasks.
func (this *TransformScheduler) Shutdown(driver scheduler.SchedulerDriver) {
	fmt.Println("Shutting down scheduler.")

	for _, taskId := range this.tasks {
		if err := this.tryKillTask(driver, taskId); err != nil {
			fmt.Printf("Failed to kill task %s\n", taskId.GetValue())
		}
	}
}

func (this *TransformScheduler) getRunningInstances() int32 {
	return atomic.LoadInt32(&this.runningInstances)
}

func (this *TransformScheduler) incRunningInstances() {
	atomic.AddInt32(&this.runningInstances, 1)
}

func (this *TransformScheduler) decRunningInstances() {
	atomic.AddInt32(&this.runningInstances, -1)
}

func (this *TransformScheduler) takePort(ports *[]*mesos.Value_Range) *uint64 {
	port := (*ports)[0].Begin
	portRange := (*ports)[0]
	portRange.Begin = proto.Uint64((*portRange.Begin) + 1)

	if *portRange.Begin > *portRange.End {
		*ports = (*ports)[1:]
	} else {
		(*ports)[0] = portRange
	}

	return port
}

func (this *TransformScheduler) createExecutor(instanceId int32, port uint64) *mesos.ExecutorInfo {
	path := strings.Split(this.config.ExecutorArchiveName, "/")
	return &mesos.ExecutorInfo{
		ExecutorId: util.NewExecutorID(fmt.Sprintf("transform-%d", instanceId)),
		Name:       proto.String("LogLine Transform Executor"),
		Source:     proto.String("cisco"),
		Command: &mesos.CommandInfo{
			Value: proto.String(fmt.Sprintf("./%s --producer.config %s --topic %s --port %d --sync %t",
				this.config.ExecutorBinaryName, this.config.ProducerConfig, this.config.Topic, port, this.config.Sync)),
			Uris: []*mesos.CommandInfo_URI{&mesos.CommandInfo_URI{
				Value:   proto.String(fmt.Sprintf("http://%s:%d/resource/%s", this.config.ArtifactServerHost, this.config.ArtifactServerPort, path[len(path)-1])),
				Extract: proto.Bool(true),
			}, &mesos.CommandInfo_URI{
				Value: proto.String(fmt.Sprintf("http://%s:%d/resource/%s", this.config.ArtifactServerHost, this.config.ArtifactServerPort, this.config.ProducerConfig)),
			}},
		},
	}
}

func (this *TransformScheduler) tryKillTask(driver scheduler.SchedulerDriver, taskId *mesos.TaskID) error {
	fmt.Printf("Trying to kill task %s\n", taskId.GetValue())

	var err error
	for i := 0; i <= this.config.KillTaskRetries; i++ {
		if _, err = driver.KillTask(taskId); err == nil {
			return nil
		}
	}
	return err
}

func (this *TransformScheduler) removeTask(id *mesos.TaskID) {
	for i, task := range this.tasks {
		if *task.Value == *id.Value {
			this.tasks = append(this.tasks[:i], this.tasks[i+1:]...)
		}
	}
}

func getScalarResources(offer *mesos.Offer, resourceName string) float64 {
	resources := 0.0
	filteredResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
		return res.GetName() == resourceName
	})
	for _, res := range filteredResources {
		resources += res.GetScalar().GetValue()
	}
	return resources
}

func getRangeResources(offer *mesos.Offer, resourceName string) []*mesos.Value_Range {
	resources := make([]*mesos.Value_Range, 0)
	filteredResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
		return res.GetName() == resourceName
	})
	for _, res := range filteredResources {
		resources = append(resources, res.GetRanges().GetRange()...)
	}
	return resources
}
