package runners

import (
	j "github.com/keenfury/axenda/job"
	"github.com/keenfury/axenda/util"
)

type API struct{}

func (a *API) WhichRunner() string {
	return "API"
}
func (a *API) RunJob(job *j.Job) error {
	hdrs := make(map[string]string, 1)
	hdrs["Content-Type"] = "application/json"
	return util.SimpleRequest("POST", job.UrlPath, &job, nil, 204, hdrs)
}
