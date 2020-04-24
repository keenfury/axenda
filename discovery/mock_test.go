package adapters

import (
	"testing"
	"time"

	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
	"github.com/stretchr/testify/assert"
)

func TestMockGetJobsSuccess(t *testing.T) {
	mock := Mock{}
	timeNow := time.Now()
	jobs, err := mock.GetJobs(timeNow)
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, 1, len(jobs), "Expected jobs count to be 1")
}

func TestMockGetJobsFailure(t *testing.T) {
	mock := Mock{}
	timeNow := time.Time{}
	_, err := mock.GetJobs(timeNow)
	assert.NotNil(t, err, "No error expected")
	assert.Equal(t, "Zero time", err.Error(), "Error should be 'Zero time'")
}

func TestMockStartJobSuccess(t *testing.T) {
	mock := Mock{Runner: &r.Mock{}}
	job := j.Job{}
	ch := make(chan j.Job)
	go func() {
		err := mock.StartJob(job, ch)
		assert.Nil(t, err, "No error expected")
	}()
	<-ch
}

func TestMockCompleteJobSuccess(t *testing.T) {
	mock := Mock{}
	job := j.Job{}
	ch := make(chan j.Job)
	go func() {
		err := mock.CompleteJob(job, ch)
		assert.Nil(t, err, "No error expected")
	}()
	<-ch
}

func TestMockWhich(t *testing.T) {
	mock := Mock{Runner: &r.Mock{}}
	msg := mock.WhichDiscovery()
	assert.Equal(t, "Mock with runner: Mock", msg)
}
