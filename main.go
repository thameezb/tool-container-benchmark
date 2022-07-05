package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"gray.net/tool-container-benchmark/benchmark"
	"gray.net/tool-container-benchmark/repo"
)

//nolint:gochecknoglobals
var runBenchmarkFunc func()

type config struct {
	logLevel         string
	dbHost           string
	dbName           string
	dbPassword       string
	dbPort           string
	dbUsername       string
	isLambda         bool
	numberOfProjects int
	numberOfEvents   int
}

func handleLambdaRequest() (events.ALBTargetGroupResponse, error) {
	runBenchmarkFunc()
	return events.ALBTargetGroupResponse{
		Body:              "benchmark running",
		StatusCode:        201,
		StatusDescription: "201 Created",
		IsBase64Encoded:   false,
		Headers:           map[string]string{},
	}, nil
}

func main() {
	cfg := getConfig()
	initLogging(cfg.logLevel)

	repo := repo.NewRepository(cfg.dbHost, cfg.dbPort, cfg.dbName, cfg.dbUsername, cfg.dbPassword)
	bm := benchmark.NewBenchmark(repo, cfg.numberOfProjects)
	runBenchmarkFunc = func() {
		bm.Run(cfg.numberOfEvents)
	}
	router := NewServer(runBenchmarkFunc)

	if !cfg.isLambda {
		log.Info("running in Server Mode")
		router.Server()
	}

	if cfg.isLambda {
		log.Info("running in Lambda Mode")
		lambda.Start(handleLambdaRequest)
	}
}

func getConfig() config {
	logLevel := getEnv("LOG_LEVEL", "info")

	dbHost := getEnv("DB_HOST")
	dbName := getEnv("DB_NAME")
	dbPassword := getEnv("DB_PASSWORD")
	dbPort := getEnv("DB_PORT")
	dbUsername := getEnv("DB_USERNAME")

	isLambda, err := strconv.ParseBool(getEnv("IS_LAMBDA", "false"))
	if err != nil {
		panic("IS_LAMBDA must be either true or false")
	}

	numberOfProjects, err := strconv.ParseInt(getEnv("NUMBER_OF_PROJECTS", "100000"), 10, 64)
	if err != nil {
		panic("NUMBER_OF_PROJECTS must be a number")
	}

	numberOfEvents, err := strconv.ParseInt(getEnv("NUMBER_OF_EVENTS", "900000"), 10, 64)
	if err != nil {
		panic("NUMBER_OF_EVENTS must be a number")
	}

	return config{
		logLevel:         logLevel,
		dbHost:           dbHost,
		dbName:           dbName,
		dbPassword:       dbPassword,
		dbPort:           dbPort,
		dbUsername:       dbUsername,
		isLambda:         isLambda,
		numberOfProjects: int(numberOfProjects),
		numberOfEvents:   int(numberOfEvents),
	}
}

func getEnv(name string, defaultValue ...string) string {
	if x, ok := os.LookupEnv(name); !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("%s is required", name))
	} else {
		return x
	}
}

func initLogging(level string) {
	l, err := log.ParseLevel(level)
	if err != nil {
		log.Fatalf("Failed to parse log level %s", level)
	}
	log.SetLevel(l)
	log.SetFormatter(&log.TextFormatter{
		DisableQuote: true,
	})

	if os.Getenv("DEVELOPMENT") != "true" {
		log.SetReportCaller(true)
	}
}
