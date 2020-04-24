package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	fr "github.com/keenfury/axenda/frequency"
	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
)

type (
	File struct {
		FileName string
		Runner   r.RunnerAdapter
	}
)

var FileRead = sync.Mutex{}

func (f *File) WhichDiscovery() string {
	return fmt.Sprintf("File with runner: %s", f.Runner.WhichRunner())
}

func (f *File) GetJobs(t time.Time) (jobs []j.Job, err error) {
	if t.IsZero() {
		err = fmt.Errorf("Zero time")
		return
	}
	// fmt.Println(t)
	newRunTime := t.Add(3 * time.Minute)
	jobsFile, errFile := f.OpenFile()
	if errFile != nil {
		err = errFile
		return
	}
	for _, j := range jobsFile {
		if newRunTime.Sub(j.RunTime) >= 0 && j.Active {
			jobs = append(jobs, j)
		}
	}
	return
}

func (f *File) StartJob(job j.Job, updateCh chan<- j.Job) (err error) {
	job.Status = "In Process"
	updateCh <- job
	return f.Runner.RunJob(&job)
}

func (f *File) CompleteJob(job j.Job, updateCh chan<- j.Job) error {
	job.Status = "Done"
	updateCh <- job
	FileRead.Lock()
	defer FileRead.Unlock()
	jobs, errFile := f.OpenFile()
	if errFile != nil {
		return errFile
	}
	for i, j := range jobs {
		if j.Token == job.Token {
			// update job frequency
			errUpdate := fr.Update(&jobs[i])
			if errUpdate != nil {
				return errUpdate
			}
		}
	}
	bJobs, errM := json.Marshal(jobs)
	if errM != nil {
		return errM
	}
	err := ioutil.WriteFile(f.FileName, bJobs, 0644)
	return err
}

func (f *File) OpenFile() (jobs []j.Job, err error) {
	bContent, errRead := ioutil.ReadFile(f.FileName)
	if errRead != nil {
		err = errRead
		return
	}
	if err = json.Unmarshal(bContent, &jobs); err != nil {
		return
	}
	return
}

/*

Sample File content format
[
	{
		"token":"unique_id_like_uuid",
		"run_time":"date_time in RFC3339 format",
		"frequency":integer see frequency.go,
		"active:true/false
	},
	{
		...
	}
]
*/
