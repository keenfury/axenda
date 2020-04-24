package runners

import (
	"fmt"

	j "github.com/keenfury/axenda/job"
)

type Mock struct{}

func (m *Mock) WhichRunner() string {
	return "Mock"
}

func (m *Mock) RunJob(job *j.Job) error {
	fmt.Println("Running this url:", job.UrlPath)
	return nil
}
