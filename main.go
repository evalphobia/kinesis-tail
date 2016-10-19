package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/kinesis"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

var (
	streamName string
	include    string
	exclude    string
	interval   int = 1
	oldest     bool
)

func init() {
	flag.StringVar(&streamName, "stream", streamName, "Kinesis Stream name to tail (*required)")
	flag.StringVar(&include, "include", include, "show the data only if this word is contained in the result")
	flag.StringVar(&exclude, "exclude", exclude, "do not show the data if this word is contained in the result")
	flag.IntVar(&interval, "interval", interval, "polling interval in sec")
	flag.BoolVar(&oldest, "oldest", oldest, "tail from oldest record (use TRIM_HORIZON)")
}

func main() {
	flag.Parse()
	err := validateFlag()
	if err != nil {
		exitWithError(err)
	}

	svc, err := kinesis.New(config.Config{})
	if err != nil {
		exitWithError(err)
	}
	svc.SetLogger(&log.StdLogger{})

	stream, err := svc.GetStream(streamName)
	if err != nil {
		exitWithError(err)
	}

	shardIDs, err := stream.GetShardIDs()
	if err != nil {
		exitWithError(err)
	}

	list := make([]*request, len(shardIDs))
	for i, id := range shardIDs {
		list[i] = &request{shardID: id}
	}

	iteratorType := kinesis.IteratorTypeLatest
	if oldest {
		iteratorType = kinesis.IteratorTypeTrimHorizon
	}

	intervalSec := time.Duration(interval) * time.Second
	for {
		for _, req := range list {
			r, err := stream.GetRecords(kinesis.GetCondition{
				ShardID:           req.shardID,
				ShardIterator:     req.shardIterator,
				ShardIteratorType: iteratorType,
			})
			if err != nil {
				fmt.Printf("[ERROR] shardID=%s; error=%s;\n", req.shardID, err.Error())
				continue
			}

			req.shardIterator = r.NextShardIterator

			printResult(r)
		}

		time.Sleep(intervalSec)
	}
}

func validateFlag() error {
	if streamName == "" {
		return fmt.Errorf("-stream is required parameter")
	}
	return nil
}

func printResult(r kinesis.RecordResult) {
	for _, rec := range r.Items {
		data := string(rec.Data)
		if include != "" && !strings.Contains(data, include) {
			continue
		}
		if exclude != "" && strings.Contains(data, exclude) {
			continue
		}
		fmt.Printf("[%s] %s\n", *rec.ApproximateArrivalTimestamp, data)
	}
}

func exitWithError(err error) {
	fmt.Printf("[ERROR] %s\n", err)
	os.Exit(2)
}

type request struct {
	shardID       string
	shardIterator string
}
