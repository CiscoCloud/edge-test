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

import ("github.com/mesos/mesos-go/scheduler"
    mesos "github.com/mesos/mesos-go/mesosproto"
    util "github.com/mesos/mesos-go/mesosutil"
    "fmt"
    "sync/atomic"
    "github.com/golang/protobuf/proto"
    "strings"
)

type TransformSchedulerConfig struct {
    // Number of CPUs allocated for each created Mesos task.
    CpuPerTask float64

    // Number of RAM allocated for each created Mesos task.
    MemPerTask float64

    // Artifact server host name. Will be used to fetch the executor.
    ArtifactServerHost  string

    // Artifact server port.Will be used to fetch the executor.
    ArtifactServerPort  int

    // Name of the executor archive file.
    ExecutorArchiveName string

    // Name of the executor binary file contained in the executor archive.
    ExecutorBinaryName  string

    // Maximum retries to kill a task.
    KillTaskRetries     int

    // Number of task instances to run.
    Instances int
}

func NewTransformSchedulerConfig() *TransformSchedulerConfig {
    return &TransformSchedulerConfig {
        CpuPerTask: 0.2,
        MemPerTask: 256,
        KillTaskRetries: 3,
    }
}

type TransformScheduler struct {
    config *TransformSchedulerConfig
    runningInstances int32
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

    offersAndTasks := make(map[*mesos.Offer][]*mesos.TaskInfo)
    for _, offer := range offers {
        cpus := getScalarResources(offer, "cpus")
        mems := getScalarResources(offer, "mem")

        fmt.Printf("Received Offer <%s> with cpus=%f, mem=%f\n", offer.Id.GetValue(), cpus, mems)

        remainingCpus := cpus
        remainingMems := mems

        var tasks []*mesos.TaskInfo
        for int(this.getRunningInstances()) < this.config.Instances && this.config.CpuPerTask <= remainingCpus && this.config.MemPerTask <= remainingMems {
            taskId := &mesos.TaskID {
                Value: proto.String(fmt.Sprintf("transform-%d", this.getRunningInstances())),
            }

            task := &mesos.TaskInfo{
                Name:     proto.String(taskId.GetValue()),
                TaskId:   taskId,
                SlaveId:  offer.SlaveId,
                Executor: this.createExecutor(this.getRunningInstances()),
                Resources: []*mesos.Resource{
                    util.NewScalarResource("cpus", float64(this.config.CpuPerTask)),
                    util.NewScalarResource("mem", float64(this.config.MemPerTask)),
                },
            }
            fmt.Printf("Prepared task: %s with offer %s for launch\n", task.GetName(), offer.Id.GetValue())

            tasks = append(tasks, task)
            remainingCpus -= this.config.CpuPerTask
            remainingMems -= this.config.MemPerTask

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
}

func (this *TransformScheduler) getRunningInstances() int32 {
    return atomic.LoadInt32(&this.runningInstances)
}

func (this *TransformScheduler) incRunningInstances() {
    atomic.AddInt32(&this.runningInstances, 1)
}

func (this *TransformScheduler) createExecutor(instanceId int32) *mesos.ExecutorInfo {
    path := strings.Split(this.config.ExecutorArchiveName, "/")
    return &mesos.ExecutorInfo{
        ExecutorId: util.NewExecutorID(fmt.Sprintf("transform-%d", instanceId)),
        Name:       proto.String("LogLine Transform Executor"),
        Source:     proto.String("cisco"),
        Command: &mesos.CommandInfo{
            Value: proto.String(fmt.Sprintf("./%s", this.config.ExecutorBinaryName)),
            Uris:  []*mesos.CommandInfo_URI{&mesos.CommandInfo_URI{
                Value: proto.String(fmt.Sprintf("http://%s:%d/%s", this.config.ArtifactServerHost, this.config.ArtifactServerPort, path[len(path)-1])),
                Extract: proto.Bool(true),
            }},
        },
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