package adapters

import (
	"fmt"
	"time"

	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
)

type (
	Mock struct {
		Runner r.RunnerAdapter
	}
)

func (m *Mock) WhichDiscovery() string {
	return fmt.Sprintf("Mock with runner: %s", m.Runner.WhichRunner())
}

func (m *Mock) GetJobs(t time.Time) (jobs []j.Job, err error) {
	fmt.Println("Mock: GetJob")
	if t.IsZero() {
		err = fmt.Errorf("Zero time")
		return
	}
	newRunTime := t.Add(1 * time.Minute)
	jobs = append(jobs, j.Job{Token: "MOCKTOKEN", RunTime: newRunTime})
	return
}

func (m *Mock) StartJob(job j.Job, updateCh chan<- j.Job) (err error) {
	fmt.Println("Mock: StartJob")
	job.Status = "In Process"
	updateCh <- job
	return m.Runner.RunJob(&job)
}

func (m *Mock) CompleteJob(job j.Job, updateCh chan<- j.Job) (err error) {
	fmt.Println("Mock: CompleteJob")
	job.Status = "Done"
	updateCh <- job
	return
}
