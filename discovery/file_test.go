package adapters

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
	"github.com/stretchr/testify/assert"
)

func TestFileGetJobsSuccess(t *testing.T) {
	fileName := "/tmp/file_test_get"
	content := []byte(`[{"token":"TOKENFILE","active":true,"run_time":"2020-01-01T00:00:00-06:00","frequency":2}]`)
	ioutil.WriteFile(fileName, content, 0644)
	defer os.Remove(fileName)
	file := File{FileName: fileName}
	timeNow := time.Now()
	jobs, err := file.GetJobs(timeNow)
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, 1, len(jobs), "Expected jobs count to be 1")
}

func TestFileGetJobsMissingFileFailure(t *testing.T) {
	file := File{}
	timeNow := time.Now()
	_, err := file.GetJobs(timeNow)
	assert.NotNil(t, err, "Error expected")
}

func TestFileGetJobsFailure(t *testing.T) {
	file := File{}
	timeNow := time.Time{}
	_, err := file.GetJobs(timeNow)
	assert.NotNil(t, err, "No error expected")
	assert.Equal(t, "Zero time", err.Error(), "Error should be 'Zero time'")
}

func TestFileStartJobSuccess(t *testing.T) {
	file := File{Runner: &r.Mock{}}
	job := j.Job{}
	ch := make(chan j.Job)
	go func() {
		err := file.StartJob(job, ch)
		assert.Nil(t, err, "No error expected")
	}()
	<-ch
}

func TestFileCompleteJobSuccess(t *testing.T) {
	fileName := "/tmp/file_test_cmp"
	content := []byte(`[{"token":"TOKENFILE","active":true,"run_time":"2020-01-01T00:00:00-06:00","frequency":2}]`)
	ioutil.WriteFile(fileName, content, 0644)
	file := File{FileName: fileName}
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	job := j.Job{Token: "TOKENFILE", RunTime: tm, Frequency: 2}
	ch := make(chan j.Job)
	go func(fileName string) {
		err := file.CompleteJob(job, ch)
		os.Remove(fileName)
		assert.Nil(t, err, "No error expected")
	}(fileName)
	<-ch
}

func TestFileCompleteJobMultipleSuccess(t *testing.T) {
	fileName := "/tmp/file_test_multiple"
	content := []byte(`[{"token":"TOKENFILE","active":true,"run_time":"2020-01-01T00:00:00-06:00","frequency":2}, {"token":"ANOTHERTOKENFILE","active":true,"run_time":"2020-01-01T00:00:00-06:00","frequency":2}]`)
	ioutil.WriteFile(fileName, content, 0644)
	file := File{FileName: fileName}
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	job := j.Job{Token: "TOKENFILE", RunTime: tm, Frequency: 2}
	ch := make(chan j.Job)
	go func(fileName string) {
		err := file.CompleteJob(job, ch)
		os.Remove(fileName)
		assert.Nil(t, err, "No error expected")
	}(fileName)
	<-ch
}

func TestFileCompleteJobFailure(t *testing.T) {
	fileName := "/tmp/file_test_fail"
	content := []byte(`[{"token":"TOKENFILE","active":true,"run_time":"2020-01-01T00:00:00-06:00","frequency":2]`)
	ioutil.WriteFile(fileName, content, 0644)
	file := File{FileName: fileName}
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	job := j.Job{Token: "TOKENFILE", RunTime: tm, Frequency: 2}
	ch := make(chan j.Job)
	go func(fileName string) {
		err := file.CompleteJob(job, ch)
		os.Remove(fileName)
		assert.NotNil(t, err, "Error expected")
		assert.Equal(t, "invalid character ']' after object key:value pair", err.Error())
	}(fileName)
	<-ch
}

func TestFileWhich(t *testing.T) {
	file := File{Runner: &r.Mock{}}
	msg := file.WhichDiscovery()
	assert.Equal(t, "File with runner: Mock", msg)
}

func TestFileCompleteJobMissingFileFailure(t *testing.T) {
	file := File{}
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	job := j.Job{Token: "TOKENFILE", RunTime: tm, Frequency: 2}
	ch := make(chan j.Job)
	go func() {
		err := file.CompleteJob(job, ch)
		assert.NotNil(t, err, "Error expected")
	}()
	<-ch
}

func TestFileCompleteJobFrequencyFailure(t *testing.T) {
	fileName := "/tmp/file_test_frequency"
	content := []byte(`[{"token":"TOKENFILE","active":true,"run_time":"2020-01-01T00:00:00-06:00","frequency":10}]`)
	ioutil.WriteFile(fileName, content, 0644)
	file := File{FileName: fileName, Runner: &r.Mock{}}
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	job := j.Job{Token: "TOKENFILE", RunTime: tm, Frequency: 10}
	ch := make(chan j.Job)
	go func(fileName string) {
		err := file.CompleteJob(job, ch)
		os.Remove(fileName)
		assert.NotNil(t, err, "Error expected")
	}(fileName)
	<-ch
}
