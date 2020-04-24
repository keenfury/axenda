package config

import (
	"os"
)

var (
	// Required: set to "true" if you want all your dates to use UTC else it will use the system's location
	// note: if using db, api or grpc, make sure your dates match in regards to the correct timezone
	UseUTC = os.Getenv("SCH_USE_UTC")
	// Optional: set to full path to the file on your local system
	JobFileName = os.Getenv("SCH_JOB_FILE_NAME")
	// Optional: set these to read from a db
	DBHost = os.Getenv("SCH_DB_HOST")
	DBUser = os.Getenv("SCH_DB_USER")
	DBPwd  = os.Getenv("SCH_DB_PWD")
	DBDB   = os.Getenv("SCH_DB_DB")
	// Optional: set these to read from an api endpoint(s)
	APIGetUrl = os.Getenv("SCH_API_GET_URL")
	APICmpUrl = os.Getenv("SCH_API_CMP_URL")
	// Optional: set this to use the grpc
	GRPCUrl = os.Getenv("SCH_GRPC_URL")
	// Optional: set to full path to push simple messages to a log file
	LogFileName = os.Getenv("SCH_LOG_FILE_NAME")
	// Optional: set either of these "true"
	UseRunnerAPI  = os.Getenv("SCH_USE_API")
	UseRunnerGRPC = os.Getenv("SCH_USE_GRPC")
)
