package frequency

import (
	"testing"
	"time"

	j "github.com/keenfury/axenda/job"
	"github.com/keenfury/axenda/util"
	"github.com/stretchr/testify/assert"
)

func TestFrequencyOnce(t *testing.T) {
	now := util.GetNow()
	job := j.Job{RunTime: now, Frequency: 1}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, false, job.Active)
}

func TestFrequencyMinute(t *testing.T) {
	now := util.GetNow()
	expecting := now.Add(1 * time.Minute)
	job := j.Job{RunTime: now, Frequency: 2}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyHourly(t *testing.T) {
	now := util.GetNow()
	expecting := now.Add(60 * time.Minute)
	job := j.Job{RunTime: now, Frequency: 3}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyDaily(t *testing.T) {
	now := util.GetNow()
	expecting := now.AddDate(0, 0, 1)
	job := j.Job{RunTime: now, Frequency: 4}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyWeekly(t *testing.T) {
	now := util.GetNow()
	expecting := now.AddDate(0, 0, 7)
	job := j.Job{RunTime: now, Frequency: 5}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyMonthly(t *testing.T) {
	now := util.GetNow()
	expecting := now.AddDate(0, 1, 0)
	job := j.Job{RunTime: now, Frequency: 6}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyQuarterly(t *testing.T) {
	now := util.GetNow()
	expecting := now.AddDate(3, 0, 0)
	job := j.Job{RunTime: now, Frequency: 7}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyMinutePast(t *testing.T) {
	now := util.GetNow()
	now = util.TruncateTimeToMinute(now)
	// move the date back a few days
	timePast := now.AddDate(0, 0, -2)
	expecting := now.Add(1 * time.Minute)
	job := j.Job{RunTime: timePast, Frequency: 2}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyHourlyPast(t *testing.T) {
	now := util.GetNow()
	now = util.TruncateTimeToMinute(now)
	// move the date back a few days
	timePast := now.AddDate(0, 0, -2)
	expecting := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), timePast.Minute(), 0, 0, now.Location()).Add(60 * time.Minute)
	job := j.Job{RunTime: timePast, Frequency: 3}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyDailyPast(t *testing.T) {
	now := util.GetNow()
	now = util.TruncateTimeToMinute(now)
	// move the date back a few days
	timePast := now.AddDate(0, 0, -2)
	expecting := now.AddDate(0, 0, 1)
	job := j.Job{RunTime: timePast, Frequency: 4}
	err := Update(&job)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, expecting, job.RunTime)
}

func TestFrequencyInvaildFrequency(t *testing.T) {
	now := util.GetNow()
	now = util.TruncateTimeToMinute(now)
	job := j.Job{RunTime: now, Frequency: 10}
	err := Update(&job)
	assert.NotNil(t, "Invalid frequency number", err.Error(), "Expecting Error")
}
