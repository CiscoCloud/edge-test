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
    "fmt"
)

type TransformScheduler struct {
    
}

func NewTransformScheduler() *TransformScheduler {
    return new(TransformScheduler)
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
    fmt.Printf("Received offers: %s\n", offers)
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