package runners

import (
	"context"
	"time"

	"github.com/keenfury/axenda/discovery/proto"
	j "github.com/keenfury/axenda/job"
	"google.golang.org/grpc"
)

type GRPC struct{}

func (g *GRPC) WhichRunner() string {
	return "GRPC"
}

func (g *GRPC) RunJob(job *j.Job) error {
	opts := grpc.WithInsecure()
	srv, errDial := grpc.Dial(job.UrlPath, opts)
	if errDial != nil {
		return errDial
	}
	defer srv.Close()
	cli := proto.NewJobServiceClient(srv)
	bPayload, errM := job.Payload.MarshalJSON()
	if errM != nil {
		return errM
	}
	pJob := &proto.Job{Token: job.Token, JobName: job.JobName, Runtime: job.RunTime.Format(time.RFC3339), Frequency: int32(job.Frequency), Payload: bPayload, Active: job.Active}
	req := proto.JobRunRequest{Job: pJob}
	_, errResp := cli.RunJob(context.Background(), &req)
	if errResp != nil {
		return errResp
	}
	return nil
}
