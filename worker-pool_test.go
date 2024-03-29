// MIT License
//
// (C) Copyright [2018, 2021] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

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
