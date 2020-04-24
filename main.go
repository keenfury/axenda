// Created by keenfury 2020
// See the LICENSE and README.md files for details

package main

import (
	"fmt"
	"time"

	"github.com/keenfury/axenda/config"
	d "github.com/keenfury/axenda/discovery"
	j "github.com/keenfury/axenda/job"
	l "github.com/keenfury/axenda/logger"
	r "github.com/keenfury/axenda/runner"
	"github.com/keenfury/axenda/util"
)

type (
	DiscoveryAdapter interface {
		WhichDiscovery() string
		GetJobs(time.Time) ([]j.Job, error)
		StartJob(j.Job, chan<- j.Job) error
		CompleteJob(j.Job, chan<- j.Job) error
	}

	LogAdapter interface {
		SetMessage(string)
	}
)

var (
	jobs             []j.Job
	JobArrayCh       chan j.Job
	JobUpdateCh      chan j.Job
	JobRemoveCh      chan j.Job
	discoveryAdapter = SetDiscoveryAdapter()
	logAdapter       = SetLoggingAdapter()
)

func main() {
	logAdapter.SetMessage(fmt.Sprintf("Using discovery: %s\n", discoveryAdapter.WhichDiscovery()))
	jobs = []j.Job{}
	JobUpdateCh = make(chan j.Job)
	minuteTicker := time.NewTicker(time.Minute)

	for {
		select {
		case job := <-JobUpdateCh:
			UpdateStatus(job, &jobs)
		case t := <-minuteTicker.C:
			noSecondsTime := util.TruncateTimeToMinute(t)
			ProcessMinute(noSecondsTime, &jobs, discoveryAdapter, JobUpdateCh)
		}
	}
}

// ProcessMinute: called by the MinuteTicker, start the process
func ProcessMinute(t time.Time, jobs *[]j.Job, ja DiscoveryAdapter, updateCh chan<- j.Job) {
	CheckForJobs(t, jobs, ja)
	RunJobs(*jobs, ja, updateCh)
}

// CheckForJobs: called by ProcessMinute, call the adpater's GetJobs, set the Job's status to 'Received'
func CheckForJobs(t time.Time, jobs *[]j.Job, ja DiscoveryAdapter) {
	newJobs, errGet := ja.GetJobs(t)
	if errGet != nil {
		logAdapter.SetMessage(fmt.Sprintf("CheckForJobs: %s", errGet))
	}
	for _, j := range newJobs {
		if CheckDup(j, *jobs) {
			j.Status = "Received"
			*jobs = append(*jobs, j)
		}
	}
}

// RunJobs: called by ProcessMinute, run Job(s) if the status has been 'Received'
// this function call the adapter's StartJob and CompleteJob
func RunJobs(jobs []j.Job, ja DiscoveryAdapter, updateCh chan<- j.Job) {
	nowWithNoSeconds := util.TruncateTimeToMinute(util.GetNow())
	for _, job := range jobs {
		if job.Status == "Received" {
			if nowWithNoSeconds.Sub(job.RunTime) >= 0 {
				go func(job j.Job) {
					if errStart := ja.StartJob(job, updateCh); errStart != nil {
						logAdapter.SetMessage(fmt.Sprintf("RunJobs: %s", errStart))
						job.Status = fmt.Sprintf("Error: %s", errStart)
						updateCh <- job
					}
					if errComplete := ja.CompleteJob(job, updateCh); errComplete != nil {
						logAdapter.SetMessage(fmt.Sprintf("CompleteJobs: %s", errComplete))
						job.Status = fmt.Sprintf("Error: %s", errComplete)
						updateCh <- job
					}
				}(job)
			}
		}
	}
}

// SetDiscoveryAdapter: look through the environment variables to deterimine with adapter to use.
// Order of precedence: file, db, api, grpc and then the failsafe mock.
func SetDiscoveryAdapter() DiscoveryAdapter {
	// check runners first so we can use them below
	var runner r.RunnerAdapter
	if config.UseRunnerAPI == "true" {
		runner = &r.API{}
	}
	if config.UseRunnerGRPC == "true" {
		runner = &r.GRPC{}
	}
	if runner == nil {
		runner = &r.Mock{}
	}
	// now check adapters
	// check for local file
	if len(config.JobFileName) > 0 {
		return &d.File{FileName: config.JobFileName, Runner: runner}
	}
	// check for DB
	if len(config.DBHost) > 0 { // will just check host
		db := d.DB{Runner: runner}
		db.Connect()
		return &db
	}
	// check for API
	if len(config.APIGetUrl) > 0 { // will just check get url
		return &d.API{Runner: runner}
	}
	// check for gRPC
	if len(config.GRPCUrl) > 0 {
		return &d.GRPC{URL: config.GRPCUrl, Runner: runner}
	}
	// default mock
	return &d.Mock{Runner: runner}
}

// SetLoggingAdapter: determines which logging adapter to use
// customize which adapter you want to use, order of precedency: file and then the failsafe stdout
func SetLoggingAdapter() LogAdapter {
	if len(config.LogFileName) > 0 {
		return &l.File{}
	}
	return &l.StdOut{}
}

// UpdateStatus: update the status of the Job in the array of Job
// remove from array of Job when status is "Done"
func UpdateStatus(job j.Job, jobs *[]j.Job) {
	removeIdx := -1
	for i, _ := range *jobs {
		if (*jobs)[i].Token == job.Token {
			(*jobs)[i].Status = job.Status
			if job.Status == "Done" {
				removeIdx = i
				break
			}
		}
	}
	if removeIdx > -1 {
		*jobs = append((*jobs)[:removeIdx], (*jobs)[removeIdx+1:]...)
	}
}

// CheckDup: checks if the Jobs already has the token, add it if needed
func CheckDup(job j.Job, jobs []j.Job) bool {
	for _, js := range jobs {
		if js.Token == job.Token {
			return false
		}
	}
	return true
}
