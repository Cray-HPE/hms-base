// Copyright (c) 2018 Cray Inc. All Rights Reserved.
package base

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

const (
	JTYPE_INVALID JobType = 0
	JTYPE_TEST    JobType = 1
	JTYPE_MAX     JobType = 2
)

var JTypeString = map[JobType]string{
	JTYPE_INVALID: "JTYPE_INVALID",
	JTYPE_TEST:    "JTYPE_TEST",
	JTYPE_MAX:     "JTYPE_MAX",
}

///////////////////////////////////////////////////////////////////////////////
// Test Helper Functions
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
// JTYPE_TEST
///////////////////////////////////////////////////////////////////////////////
type JobTest struct {
	Status JobStatus
	Num    int
	Msg    string
	Err    error
	Logger *log.Logger
}

func NewJobTest(num int, msg string, status JobStatus, lg *log.Logger) Job {
	j := new(JobTest)
	j.Status = status
	j.Num = num
	j.Msg = msg
	j.Logger = lg

	if lg == nil {
		j.Logger = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	} else {
		j.Logger = lg
	}
	return j
}

// Log to logging infrastructure.
func (j *JobTest) Log(format string, a ...interface{}) {
	// Use caller's line number (depth=2)
	j.Logger.Output(2, fmt.Sprintf(format, a...))
}

func (j *JobTest) Type() JobType {
	return JTYPE_TEST
}

func (j *JobTest) Run() {
	//Do stuff what is here is temporary
	//time.Sleep(5 * time.Second)
	j.Log("Processing Test Job #%d with Status='%s' Msg='%s'", j.Num, JStatString[j.Status], j.Msg)
}

func (j *JobTest) GetStatus() (JobStatus, error) {

	if j.Status == JSTAT_ERROR {
		return j.Status, j.Err
	}
	return j.Status, nil
}

func (j *JobTest) SetStatus(newStatus JobStatus, err error) (JobStatus, error) {
	if newStatus >= JSTAT_MAX {
		return j.Status, fmt.Errorf("Error: Invalid Status")
	} else {
		oldStatus := j.Status
		j.Status = newStatus
		j.Err = err
		return oldStatus, nil
	}
}

// This JobType does not support cancelling the job while it is being processed
func (j *JobTest) Cancel() JobStatus {
	if j.Status == JSTAT_QUEUED || j.Status == JSTAT_DEFAULT {
		j.Status = JSTAT_CANCELLED
	}
	return j.Status
}

///////////////////////////////////////////////////////////////////////////////
// Pre-Test Setup
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
// Unit Tests
///////////////////////////////////////////////////////////////////////////////

// Tests the worker pool by making 10 workers and 10 jobs
func TestJobs(t *testing.T) {
	wp := NewWorkerPool(10, 10)
	wp.Run()
	jobList := make([]Job, 10)
	for i, _ := range jobList {
		jobList[i] = NewJobTest(i, "I'm a test Job", JSTAT_DEFAULT, nil)
		wp.Queue(jobList[i])
	}
	// Wait for all jobs to complete
	for _, job := range jobList {
		for {
			if status, _ := job.GetStatus(); status == JSTAT_COMPLETE {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}
	wp.Stop()
}
