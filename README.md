# Axenda
Axenda => Galician for 'schedule' (it sounds cool).  A generic scheduler written in golang.

## Overview

Why a scheduler in its own application?  I have written my share of API(s) in Golang.  Yes, you can shove some type of go routine in your API code, off of the main process thread and call it good.  To me that feels like a hack.  It the world of microservices, you have a small code base that runs one job.

I really like the idea of microservices.  This scheduler is a microservice that I have written a few times, though it was written to the environment I was working in at the time.  I got a wild idea to make it more dynamic and user friendly, so here it is.

I've written this scheduler to do one thing, run every minute, looking for jobs through a discovery adapter to run and then send a 'message' through a runner adapter to somewhere else to run that job.  That is it.  All scheduler is worried about is the jobs it knows about.

The way I've written this scheduler, I think, is more complete and would help anyone get up and running very quickly.

I've built the scheduler to be dynamic in where it can find jobs.  The list of discovery adapters are:

- File
- Database table
- API service
- GRPC service

Depending of the environment variables you set, scheduler will use the appropriate 'discovery' method.

When a job runs on the correct RunTime then a 'message' is sent through the runner adapter, two known runner adapters:

- API service
- GRPC service

The scheduler doesn't have any code to run any of the jobs, that would up for you to decide where that code lies; in some type of endpoint on an API service or a url of a GRPC service.  I guess you could develop another way and add it as an 'runner'.  The scheduler only cares about which jobs are due to run and where to send the 'message'.

## Job
A job is the heart of the scheduler and its data structure has some key fields.

- Token: [string] unique (uuid for example)
- JobName: [string]
- RunTime: [datetime] RFC3339 format
- Frequency: [integer]
- Active: [boolean]
- Payload: [bytes]
- Status: [string] used only within the app

## Discovery
Where does scheduler find its jobs?  These are explained below and are in order of presedence.

### File
If the environment variable of SCH_JOB_FILE_NAME is set then scheduler will look at the full path saved in this environment variable.

e.g. export SCH_JOB_FILE_NAME=/path/to/your/file

Format of this file can be found in discovery/file.go.  FYI: you can set the file content to "pretty" json but once to program runs and saves back out the content, the "pretty" format will go away.

### Database
If the environment variable of SCH_DB_HOST is set then the scheduler will look at a database table for jobs.

See discovery/db.go for information on table schema.

Other environment variables that may need to be set:

- SCH_DB_USER
- SCH_DB_PWD
- SCH_DB_DB

### API
If the environment variable of SCH_API_GET_URL is set then the scheduler will look for the jobs through an API endpoint.

The SCH_API_GET_URL will add an Unix date as part of the url route.

SCH_API_CMP_URL endpoint will be a POST request with the job in the body.

All data transfer format will be in JSON.

### GRPC
If the environment variable of SCH_GRPC_URL is set then the scheduler will look for the jobs through a GRPC endpoint.

All data transfer will be using the protobuf, see discovery/proto/grpc.proto

## Runner

### API
The easiest way to run your job is to have a dedicated endpoint route that will verify the token being sent by the scheduler.  That code will do whatever you need.  The url should be saved per job in the url_path.

### GDPR
The way this protobuf code is written there is one server endpoint called 'RunJob' that will take any job call from this service.  This makes it a bit tricky to decide on your side what code to call.  I'll leave that up to you, maybe with some modifications on the table side you can have some branching code.  I like the API way because the router will do that for you.

## Other Features

### Frequency
So in order to run a job multiple times you will need to figure out when the next time will be.  The scheduler is designed take the RunTime and increment the time based on the frequency the job is set at.

See frequency/frequency.go for the list of options.

Note: in order to process date/time(s) in the past, the date/time needs to be calculated to current date but keep the RunTime's hour and minute with the exception of:

- Minute frequency: past date/time => now
- Hour frequency: past date/time => now with RunTime's minute

### Logging
I've also include an easy way to direct logging to either:

- STDOUT
- File

Though it would take too much to send logging to an API endpoint or something custom, like a messaging queue.  (See main.go function SetLoggingAdapter)