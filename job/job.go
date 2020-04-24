package jobs

import (
	"encoding/json"
	"time"
)

type (
	Job struct {
		Token     string          `db:"token" json:"token"`
		JobName   string          `db:"job_name" json:"job_name"`
		RunTime   time.Time       `db:"run_time" json:"run_time"`
		UrlPath   string          `db:"url_path" json:"url_path"`
		Frequency int             `db:"frequency" json:"frequency"`
		Active    bool            `db:"active" json:"active"`
		Payload   json.RawMessage `db:"payload" json:"payload"`
		Status    string          `json:"-"`
	}
)
