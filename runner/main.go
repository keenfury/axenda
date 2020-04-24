package runners

import (
	j "github.com/keenfury/axenda/job"
)

type RunnerAdapter interface {
	WhichRunner() string
	RunJob(job *j.Job) error
}
