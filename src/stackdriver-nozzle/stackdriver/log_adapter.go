package stackdriver

import (
	"time"

	"cloud.google.com/go/logging"
	"github.com/cloudfoundry/lager"
	"golang.org/x/net/context"
)

const (
	logId                = "cf_logs"
	DefaultBatchCount    = "10"
	DefaultBatchDuration = "1s"
)

type LogAdapter interface {
	PostLog(*Log)
}

type Log struct {
	Payload interface{}
	Labels  map[string]string
}

func NewLogAdapter(projectID string, batchCount int, batchDuration time.Duration, logger lager.Logger) (LogAdapter, error) {
	loggingClient, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}

	loggingClient.OnError = func(err error) {
		logger.Fatal("stackdriverClientOnError", err)
	}

	sdLogger := loggingClient.Logger(logId,
		logging.EntryCountThreshold(batchCount),
		logging.DelayThreshold(batchDuration),
	)

	return &logClient{
		sdLogger:  sdLogger,
		projectID: projectID,
		logger:    logger,
	}, nil
}

type logClient struct {
	sdLogger  *logging.Logger
	projectID string
	logger    lager.Logger
}

func (s *logClient) PostLog(log *Log) {
	entry := logging.Entry{
		Payload: log.Payload,
		Labels:  log.Labels,
	}
	s.sdLogger.Log(entry)
}
