package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/keenfury/axenda/discovery/proto"
	f "github.com/keenfury/axenda/frequency"
	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
	"google.golang.org/grpc"
)

type (
	GRPC struct {
		URL    string
		Runner r.RunnerAdapter
	}
)

func (g *GRPC) WhichDiscovery() string {
	return fmt.Sprintf("GRPC with runner: %s", g.Runner.WhichRunner())
}

func (g *GRPC) GetJobs(t time.Time) (jobs []j.Job, err error) {
	newRunTime := t.Add(3 * time.Minute)
	if t.IsZero() {
		err = fmt.Errorf("Zero time")
		return
	}
	jg := proto.JobGetRequest{Runtime: newRunTime.Format(time.RFC3339)}
	opts := grpc.WithInsecure()
	srv, errDial := grpc.Dial(g.URL, opts)
	if errDial != nil {
		err = errDial
		return
	}
	defer srv.Close()
	cli := proto.NewJobServiceClient(srv)
	resp, errResp := cli.GetJob(context.Background(), &jg)
	if errResp != nil {
		err = errResp
		return
	}
	for _, r := range resp.Jobs {
		timeParse, errParse := time.Parse(time.RFC3339, r.Runtime)
		if errParse != nil {
			fmt.Println("error in parsing time") // TODO: log
			continue
		}
		jobs = append(jobs, j.Job{Token: r.Token, RunTime: timeParse, Frequency: int(r.Frequency)})
	}
	return
}

func (g *GRPC) StartJob(job j.Job, updateCh chan<- j.Job) (err error) {
	job.Status = "In Process"
	updateCh <- job
	return g.Runner.RunJob(&job)
}

func (g *GRPC) CompleteJob(job j.Job, updateCh chan<- j.Job) error {
	job.Status = "Done"
	updateCh <- job
	errUpdate := f.Update(&job)
	if errUpdate != nil {
		return errUpdate
	}
	// set up GRPC job struct
	runTimeStr := job.RunTime.Format(time.RFC3339)
	pj := proto.Job{Token: job.Token, Runtime: runTimeStr, Frequency: int32(job.Frequency), Active: job.Active}
	opts := grpc.WithInsecure()
	srv, errDial := grpc.Dial(g.URL, opts)
	if errDial != nil {
		return errDial
	}
	defer srv.Close()
	cli := proto.NewJobServiceClient(srv)
	req := proto.JobCmpRequest{Job: &pj}
	_, errResp := cli.CmpJob(context.Background(), &req)
	if errResp != nil {
		return errResp
	}
	return nil
}
