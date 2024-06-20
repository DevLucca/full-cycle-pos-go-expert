package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	}()

	endpoint, _ := redisC.Endpoint(ctx, "")
	os.Setenv("REDIS_ADDR", endpoint)
	os.Setenv("RATE_LIMITER_TIMEFRAME", "RPM")
	os.Setenv("RATE_LIMITER_MAX_REQUESTS_PER_TIME", "5")
	os.Setenv("RATE_LIMITER_HEADER_LIMITER", "{\"x-api-key\": 1}")

	go main()

	os.Exit(m.Run())
}

func Test_IPRPMRate(t *testing.T) {
	requestCount := 100

	wg := sync.WaitGroup{}
	wg.Add(requestCount)

	var successCount atomic.Int64

	for range requestCount {
		go func() {
			defer wg.Done()

			r, _ := http.NewRequest("GET", "http://localhost:8080", nil)
			r.Header.Add("X-Real-Ip", "ip-test")

			res, _ := http.DefaultClient.Do(r)

			if res.StatusCode == http.StatusOK {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(5), successCount.Load())
}

func Test_NotEqualIPRPMRate(t *testing.T) {
	requestCount := 100

	wg := sync.WaitGroup{}
	wg.Add(requestCount)

	var successCount atomic.Int64

	for i := range requestCount {
		go func() {
			defer wg.Done()

			r, _ := http.NewRequest("GET", "http://localhost:8080", nil)
			r.Header.Add("X-Real-Ip", fmt.Sprintf("ip-test-%d", i))

			res, _ := http.DefaultClient.Do(r)

			if res.StatusCode == http.StatusOK {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(requestCount), successCount.Load())
}

func Test_HeaderRPMRate(t *testing.T) {
	requestCount := 100

	wg := sync.WaitGroup{}
	wg.Add(requestCount)

	var successCount atomic.Int64

	for range requestCount {
		go func() {
			defer wg.Done()

			r, _ := http.NewRequest("GET", "http://localhost:8080", nil)
			r.Header.Add("x-api-key", "api-key-test")

			res, _ := http.DefaultClient.Do(r)

			if res.StatusCode == http.StatusOK {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(1), successCount.Load())
}

func Test_NotEqualHeaderRPMRate(t *testing.T) {
	requestCount := 100

	wg := sync.WaitGroup{}
	wg.Add(requestCount)

	var successCount atomic.Int64

	for i := range requestCount {
		go func() {
			defer wg.Done()

			r, _ := http.NewRequest("GET", "http://localhost:8080", nil)
			r.Header.Add("x-api-key", fmt.Sprintf("api-key-test-%d", i))

			res, _ := http.DefaultClient.Do(r)

			if res.StatusCode == http.StatusOK {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(requestCount), successCount.Load())
}

func Test_RPMRate(t *testing.T) {
	requestCount := 100

	wg := sync.WaitGroup{}
	wg.Add(requestCount)

	var successCount atomic.Int64

	for i := range requestCount {
		go func() {
			defer wg.Done()

			r, _ := http.NewRequest("GET", "http://localhost:8080", nil)

			if i%2 == 0 {
				r.Header.Add("X-Real-Ip", "ip-test")
			} else {
				r.Header.Add("x-api-key", "api-key")
			}

			res, _ := http.DefaultClient.Do(r)

			if res.StatusCode == http.StatusOK {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(6), successCount.Load())
}
