package adapters

import (
	"fmt"
	"time"

	"github.com/keenfury/axenda/config"
	f "github.com/keenfury/axenda/frequency"
	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
	"github.com/keenfury/axenda/util"
)

type (
	API struct {
		Runner r.RunnerAdapter
	}
)

func (a *API) WhichDiscovery() string {
	return fmt.Sprintf("API with runner: %s", a.Runner.WhichRunner())
}

func (a *API) GetJobs(t time.Time) (jobs []j.Job, err error) {
	if t.IsZero() {
		err = fmt.Errorf("Zero time")
		return
	}
	newRunTime := t.Add(3 * time.Minute)
	url := fmt.Sprintf("%s/%d", config.APIGetUrl, newRunTime.Unix())
	err = util.SimpleRequest("GET", url, nil, &jobs, 200, nil)
	return
}

func (a *API) StartJob(job j.Job, updateCh chan<- j.Job) (err error) {
	job.Status = "In Process"
	updateCh <- job
	a.Runner.RunJob(&job)
	return
}

func (a *API) CompleteJob(job j.Job, updateCh chan<- j.Job) (err error) {
	job.Status = "Done"
	updateCh <- job
	errUpdate := f.Update(&job)
	if errUpdate != nil {
		return errUpdate
	}
	url := fmt.Sprintf("%s", config.APICmpUrl)
	hdrs := make(map[string]string, 1)
	hdrs["Content-Type"] = "application/json"
	return util.SimpleRequest("POST", url, &job, nil, 200, hdrs)
}
