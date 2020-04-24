package adapters

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/keenfury/axenda/config"
	f "github.com/keenfury/axenda/frequency"
	j "github.com/keenfury/axenda/job"
	r "github.com/keenfury/axenda/runner"
	_ "github.com/lib/pq"
)

/*
This adapter will directly talk to your DB and interact with the table with the given config parameter set correctly
- DBHost
- DBUser
- DBPwd
- DBDB

DB adapter is built for Postgres but can be changed to run vs. another DB engine depending on your situation, you will need to change the follow:
- library import of the DB engine
- connectionStr to fit the correct format
- connection with correct db engine
- table syntax

see helper syntax at the end of this file
*/

type (
	DB struct {
		DB     *sqlx.DB
		Runner r.RunnerAdapter
	}
)

func (d *DB) Connect() (err error) {
	connectionStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", config.DBUser, config.DBPwd, config.DBDB, config.DBHost)

	d.DB, err = sqlx.Connect("postgres", connectionStr)
	if err != nil {
		log.Panicln("Could not connect with connection string:", connectionStr) // TODO: log and determine if we panic
	}
	return
}

func (d *DB) WhichDiscovery() string {
	return fmt.Sprintf("DB with runner: %s", d.Runner.WhichRunner())
}

func (d *DB) GetJobs(t time.Time) (jobs []j.Job, err error) {
	if t.IsZero() {
		err = fmt.Errorf("Zero time")
		return
	}
	sqlSelect := "select token, run_time, url_path, frequency from schedule where run_time < $1 and active = true"
	runTimeEnd := t.Add(3 * time.Minute)
	jobsDB := []j.Job{}
	errSelect := d.DB.Select(&jobsDB, sqlSelect, runTimeEnd)
	if errSelect != nil {
		return nil, errSelect
	} else {
		jobs = append(jobs, jobsDB...)
	}
	return
}

func (d *DB) StartJob(job j.Job, updateCh chan<- j.Job) error {
	job.Status = "In Process"
	updateCh <- job
	return d.Runner.RunJob(&job)
}

func (d *DB) CompleteJob(job j.Job, updateCh chan<- j.Job) (err error) {
	job.Status = "Done"
	updateCh <- job
	errUpdate := f.Update(&job)
	if errUpdate != nil {
		return errUpdate
	}
	sqlUpdate := "update schedule set run_time = $1 where token = $2"
	_, errExec := d.DB.Exec(sqlUpdate, job.RunTime, job.Token)
	if errExec != nil {
		return errExec
	}
	return
}

/*
This table syntax will help you set a table for this adapter to be used correctly (though you can change what you want, you will need to change the
struct in job.go)
create table schedule (
	token uuid not null primary key,
	job_name string null,
	run_time timestamp not null,
	url_path text not null,
	frequency int not null,
	payload json null,
	active boolean not null default true
);

- token (uuid): unique identifier (I like it versus an int just in case, harder to guess)
- job_name (string): name for the job
- run_time (timestamp/datetime/etc): the date/time you want this job to run down to the minute
- url_path (string): the path whether an api endpoint path or gRPC path, basically where the job is going to run
- frequency (int): [optional] determine how to move the date to the next run_time (see table below)
- payload (json): [optional] although this program doesn't use it (though you could add it), what every program runs your task you can read the
	payload optionally and store all kinds of info in there for your task)
- active (bool): self-explanatory

- the following should follow the list found in frequency.go
create table frequency (
	id serial primary key,
	frequency_name varchar(20)
);

insert into frequency (id, frequency_name) values (1, 'Once');
insert into frequency (id, frequency_name) values (2, 'Minute');
insert into frequency (id, frequency_name) values (3, 'Hour');
insert into frequency (id, frequency_name) values (4, 'Daily');
insert into frequency (id, frequency_name) values (5, 'Weekly');
insert into frequency (id, frequency_name) values (6, 'Monthly');
insert into frequency (id, frequency_name) values (7, 'Quarterly');
*/
