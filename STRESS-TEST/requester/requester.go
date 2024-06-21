package requester

import (
	"context"
	"log"
	"net/http"
	"os"
	"slices"
	"sync"
	"text/template"
	"time"

	_ "embed"
)

type response struct {
	StartedAt time.Time
	EndedAt   time.Time

	*http.Response
	Error error
}

func Request(ctx context.Context, url string, concurrency, requests int) {
	startTime := time.Now()
	totalRequestCount := requests
	responses := make(chan response, requests)

	var dispatched int

	wg := sync.WaitGroup{}
	waiter := make(chan struct{})
	for totalRequestCount > 0 {
		go func() {
			defer wg.Done()

			<-waiter
			startedAt := time.Now()
			resp, err := http.Get(url)
			endedAt := time.Now()

			responses <- response{
				StartedAt: startedAt,
				EndedAt:   endedAt,
				Response:  resp,
				Error:     err,
			}
		}()
		dispatched++
		totalRequestCount--

		if (concurrency%dispatched == 0 && dispatched == concurrency) ||
			(totalRequestCount == 0) {
			wg.Add(dispatched)
			close(waiter)
			wg.Wait()
			waiter = make(chan struct{})
			dispatched = 0
		}

	}

	close(responses)
	report(responses, requests, time.Since(startTime))
}

//go:embed report.tmpl
var tmplReport []byte

func report(responses <-chan response, requestCount int, doneIn time.Duration) {

	var errors []error
	var statusCodeDistributionCount = map[int]int{}
	var responseTimes []time.Duration
	for res := range responses {
		if res.Error != nil {
			errors = append(errors, res.Error)
			continue
		}
		statusCodeClass := res.StatusCode / 100
		statusCodeDistributionCount[statusCodeClass]++

		responseTimes = append(responseTimes, res.EndedAt.Sub(res.StartedAt))
	}

	tmpl, _ := template.New("report").Parse(string(tmplReport))

	err := tmpl.Execute(os.Stdout, map[string]any{
		"totalRequests":               requestCount,
		"completedRequests":           requestCount - len(errors),
		"statusCodeDistributionCount": statusCodeDistributionCount,
		"totalTime":                   doneIn.String(),
		"p99th":                       getLatency(responseTimes, 99),
		"p75th":                       getLatency(responseTimes, 75),
		"p50th":                       getLatency(responseTimes, 50),
		"errorCount":                  len(errors),
		"errors":                      errors,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func getLatency(list []time.Duration, percentile int) time.Duration {
	slices.Sort(list)

	index := int(float64(len(list)) * (float64(percentile) / 100))
	return list[index]
}
