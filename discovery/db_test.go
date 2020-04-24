package adapters

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
	"github.com/stretchr/testify/assert"
)

func TestDBGetJobsSuccess(t *testing.T) {
	// set up mock db
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	rows := mock.NewRows([]string{"token", "run_time", "url_path", "frequency"}).AddRow("TOKENDB", tm, "", 4)
	mock.ExpectQuery("select (.+) from schedule where run_time").WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := DB{DB: sqlxDB}
	timeNow := time.Now()
	jobs, err := db.GetJobs(timeNow)
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, 1, len(jobs), "Expected jobs count to be 1")
}

func TestDBGetJobsFailure(t *testing.T) {
	db := DB{}
	timeNow := time.Time{}
	_, err := db.GetJobs(timeNow)
	assert.NotNil(t, err, "Error expected")
	assert.Equal(t, "Zero time", err.Error(), "Error should be 'Zero time'")
}

func TestDBGetJobsDBFailure(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	mock.ExpectQuery("select (.+) from schedule where run_time").WillReturnError(fmt.Errorf("select error"))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := DB{DB: sqlxDB}
	timeNow := time.Now()
	_, err := db.GetJobs(timeNow)
	assert.NotNil(t, err, "Error expected")
	assert.Equal(t, "select error", err.Error(), "Error should be 'select error'")
}
func TestDBStartJobSuccess(t *testing.T) {
	db := DB{Runner: &r.Mock{}}
	job := j.Job{}
	ch := make(chan j.Job)
	go func() {
		err := db.StartJob(job, ch)
		assert.Nil(t, err, "No error expected")
	}()
	<-ch
}

func TestDBCompleteJobSuccess(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	mockDB, mock, _ := sqlmock.New()
	mock.ExpectExec("update schedule").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := DB{DB: sqlxDB}
	job := j.Job{RunTime: tm, Frequency: 4}
	ch := make(chan j.Job)
	go func() {
		err := db.CompleteJob(job, ch)
		assert.Nil(t, err, "No error expected")
		mockDB.Close()
	}()
	<-ch
}

func TestDBCompleteJobDBError(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	mockDB, mock, _ := sqlmock.New()
	mock.ExpectExec("update schedule").WillReturnError(fmt.Errorf("update error"))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := DB{DB: sqlxDB}
	job := j.Job{RunTime: tm, Frequency: 4}
	ch := make(chan j.Job)
	go func() {
		err := db.CompleteJob(job, ch)
		assert.NotNil(t, err, "Error expected")
		assert.Equal(t, "update error", err.Error(), "should give error of 'update error'")
		mockDB.Close()
	}()
	<-ch
}

func TestDBCompleteJobUpdateError(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00-06:00")
	db := DB{}
	job := j.Job{RunTime: tm, Frequency: 10}
	ch := make(chan j.Job)
	go func() {
		err := db.CompleteJob(job, ch)
		assert.NotNil(t, err, "Error expected")
		assert.Equal(t, "Invalid frequency number", err.Error(), "should give error of 'Invalid frequency number'")
	}()
	<-ch
}
func TestDBWhich(t *testing.T) {
	db := DB{Runner: &r.Mock{}}
	msg := db.WhichDiscovery()
	assert.Equal(t, "DB with runner: Mock", msg)
}
