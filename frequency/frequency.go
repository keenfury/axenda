package frequency

import (
	"fmt"
	"time"

	j "github.com/keenfury/axenda/job"
	"github.com/keenfury/axenda/util"
)

const (
	Once = 1 + iota
	Minute
	Hourly
	Daily
	Weekly
	Monthly
	Quarterly
)

func Update(job *j.Job) (err error) {
	// if the time is in the past, bring it to the current date, keeping the same hour/minute
	truncNow := util.TruncateTimeToMinute(util.GetNow())
	truncTime := util.TruncateTimeToMinute(job.RunTime)
	if truncNow.Sub(truncTime) > 0 {
		if job.Frequency == Minute {
			job.RunTime = time.Date(truncNow.Year(), truncNow.Month(), truncNow.Day(), truncNow.Hour(), truncNow.Minute(), 0, 0, truncTime.Location())
		} else if job.Frequency == Hourly {
			job.RunTime = time.Date(truncNow.Year(), truncNow.Month(), truncNow.Day(), truncNow.Hour(), truncTime.Minute(), 0, 0, truncTime.Location())
		} else {
			job.RunTime = time.Date(truncNow.Year(), truncNow.Month(), truncNow.Day(), truncTime.Hour(), truncTime.Minute(), 0, 0, truncTime.Location())
		}
	}
	switch int(job.Frequency) {
	case Once:
		job.Active = false
	case Minute:
		job.RunTime = job.RunTime.Add(1 * time.Minute)
	case Hourly:
		job.RunTime = job.RunTime.Add(60 * time.Minute)
	case Daily:
		job.RunTime = job.RunTime.AddDate(0, 0, 1)
	case Weekly:
		job.RunTime = job.RunTime.AddDate(0, 0, 7)
	case Monthly:
		job.RunTime = job.RunTime.AddDate(0, 1, 0)
	case Quarterly:
		job.RunTime = job.RunTime.AddDate(3, 0, 0)
	default:
		err = fmt.Errorf("Invalid frequency number")
	}
	return
}
