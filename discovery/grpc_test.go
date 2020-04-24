package adapters

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	p "github.com/keenfury/axenda/discovery/proto"
	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
)

type Server struct{}

func (s *Server) GetJob(ctx context.Context, gr *p.JobGetRequest) (*p.JobGetResponse, error) {
	jb := p.JobGetResponse{}
	if gr.Runtime == "1972-01-25T00:03:00-06:00" {
		// simulate error
		return &jb, fmt.Errorf("Error from server")
	}
	job := &p.Job{Token: "FROMRPC", Runtime: "2020-04-13T15:00:00-06:00", Frequency: 3}
	if gr.Runtime == "1972-01-25T10:03:00-06:00" {
		// simulate error
		job.Runtime = "2020-04-13T15:00:00-06:0"
	}
	jb.Jobs = append(jb.Jobs, job)
	return &jb, nil
}

func (s *Server) CmpJob(ctx context.Context, req *p.JobCmpRequest) (*p.JobCmpResponse, error) {
	if req.Job.Token == "ERROR_TOKEN" {
		return nil, fmt.Errorf("Error here")
	}
	m := p.JobCmpResponse{}
	m.Message = "All updated"
	return &m, nil
}

func (s *Server) RunJob(ctx context.Context, req *p.JobRunRequest) (*p.JobRunResponse, error) {
	if req.Job.Token == "ERROR_TOKEN" {
		return nil, fmt.Errorf("Error here")
	}
	m := p.JobRunResponse{}
	m.Message = "All updated"
	return &m, nil
}

func Serve() {
	lis, err := net.Listen("tcp", ":12500")
	if err != nil {
		fmt.Println("Unable to listen on port 12500")
		panic("Doh")
	}
	s := grpc.NewServer()
	p.RegisterJobServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		fmt.Println("Unable to serve")
		panic("Doh")
	}
}

func TestMain(m *testing.M) {
	go Serve()
	os.Exit(m.Run())
}

func TestGRPCGetJobsSuccess(t *testing.T) {
	rpc := GRPC{URL: "localhost:12500"}
	timeNow := time.Now()
	jobs, err := rpc.GetJobs(timeNow)
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, 1, len(jobs), "Expected jobs count to be 1")
}

func TestGRPCGetJobsZeroTimeFailure(t *testing.T) {
	rpc := GRPC{URL: "localhost:12500"}
	timeNow := time.Time{}
	_, err := rpc.GetJobs(timeNow)
	assert.NotNil(t, err, "No error expected")
	assert.Equal(t, "Zero time", err.Error(), "Error should be 'Zero time'")
}

func TestGRPCGetJobsDailFailure(t *testing.T) {
	rpc := GRPC{URL: "localhost:12501"}
	timeNow := time.Now()
	_, err := rpc.GetJobs(timeNow)
	assert.NotNil(t, err, "Error expected")
}

func TestGRPCGetJobsServerFailure(t *testing.T) {
	rpc := GRPC{URL: "localhost:12500"}
	timeNow, _ := time.Parse(time.RFC3339, "1972-01-25T00:00:00-06:00")
	_, err := rpc.GetJobs(timeNow)
	assert.NotNil(t, err, "Error expected")
	assert.Equal(t, "rpc error: code = Unknown desc = Error from server", err.Error(), "Error should be 'rpc error: code = Unknown desc = Error from server'")
}

func TestGRPCGetJobsParseTimeFailure(t *testing.T) {
	rpc := GRPC{URL: "localhost:12500", Runner: &r.Mock{}}
	timeNow, _ := time.Parse(time.RFC3339, "1972-01-25T10:00:00-06:00")
	jobs, err := rpc.GetJobs(timeNow)
	assert.Nil(t, err, "Error expected")
	assert.Equal(t, 0, len(jobs), "No jobs made it")
}

func TestGRPCStartJobSuccess(t *testing.T) {
	rpc := GRPC{URL: "localhost:12500", Runner: &r.Mock{}}
	job := j.Job{}
	ch := make(chan j.Job)
	go func() {
		err := rpc.StartJob(job, ch)
		assert.Nil(t, err, "No error expected")
	}()
	<-ch
}

func TestGRPCStartJobServerFailure(t *testing.T) {
	rpc := GRPC{URL: "localhost:12501", Runner: &r.GRPC{}}
	job := j.Job{Token: "ERROR_TOKEN", Frequency: 4}
	ch := make(chan j.Job)
	go func() {
		err := rpc.StartJob(job, ch)
		assert.NotNil(t, err, "Error expected")
	}()
	<-ch
}

func TestGRPCCompleteJobSuccess(t *testing.T) {
	rpc := GRPC{URL: "localhost:12500"}
	job := j.Job{Frequency: 4}
	ch := make(chan j.Job)
	go func() {
		err := rpc.CompleteJob(job, ch)
		assert.Nil(t, err, "No error expected")
	}()
	<-ch
}

func TestGRPCCompleteJobDailFailure(t *testing.T) {
	rpc := GRPC{URL: "localhost:12501"}
	job := j.Job{Frequency: 4}
	ch := make(chan j.Job)
	go func() {
		err := rpc.CompleteJob(job, ch)
		assert.NotNil(t, err, "No error expected")
	}()
	<-ch
}

func TestGRPCD(t *testing.T) {
	rpc := GRPC{Runner: &r.Mock{}}
	msg := rpc.WhichDiscovery()
	assert.Equal(t, "GRPC with runner: Mock", msg)
}
