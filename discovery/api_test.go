package adapters

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/keenfury/axenda/config"
	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
	"github.com/stretchr/testify/assert"
)

func TestAPIGetJobsSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `[{"token":"MOCKAPI","run_time": "2020-04-13T15:00:00-06:00","frequency": 3}]`)
	}))
	defer srv.Close()
	config.APIGetUrl = srv.URL
	api := API{}
	timeNow := time.Now()
	jobs, err := api.GetJobs(timeNow)
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, 1, len(jobs), "Expected jobs count to be 1")
}

func TestAPIGetJobsFailure(t *testing.T) {
	api := API{}
	timeNow := time.Time{}
	_, err := api.GetJobs(timeNow)
	assert.NotNil(t, err, "No error expected")
	assert.Equal(t, "Zero time", err.Error(), "Error should be 'Zero time'")
}

func TestAPIStartJobSuccess(t *testing.T) {
	api := API{Runner: &r.Mock{}}
	job := j.Job{}
	ch := make(chan j.Job)
	go func() {
		err := api.StartJob(job, ch)
		assert.Nil(t, err, "No error expected")
	}()
	<-ch
}

func TestAPICompleteJobSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"all good"}`)
	}))
	config.APICmpUrl = srv.URL
	api := API{}
	job := j.Job{Token: "TOKENAPI", Frequency: 4, RunTime: time.Now()}
	ch := make(chan j.Job)
	go func() {
		err := api.CompleteJob(job, ch)
		assert.Nil(t, err, "No error expected")
		srv.Close()
	}()
	<-ch
}

func TestAPIWhich(t *testing.T) {
	api := API{Runner: &r.Mock{}}
	msg := api.WhichDiscovery()
	assert.Equal(t, "API with runner: Mock", msg)
}
