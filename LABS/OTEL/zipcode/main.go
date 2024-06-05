package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/OTEL/zipcode/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	WEATHER_API = "http://weather-api:8001"
)

func init() {
	os.Setenv("APP_NAME", "zipcode-api")
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdown, err := telemetry.SetupProvider(ctx, os.Getenv("APP_NAME"))
	if err != nil {
		return
	}

	mux := http.NewServeMux()
	t := otel.Tracer("zipcode")

	// Register handlers.
	mux.HandleFunc("POST /{zipcode}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := t.Start(ctx, "zipcode")
		defer span.End()

		zipcode := r.PathValue("zipcode")

		if matched, _ := regexp.Match(`^\d{8}$`, []byte(zipcode)); !matched {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", WEATHER_API, zipcode), nil)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data any

		if 300 <= res.StatusCode {
			body, _ := io.ReadAll(res.Body)
			http.Error(w, string(body), res.StatusCode)
			return
		}

		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if code := res.StatusCode; code != 200 {
			http.Error(w, "status code not expected", code)
			return
		}

		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	srv := &http.Server{
		Addr:         ":8000",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
	defer func() {
		err = errors.Join(err, shutdown(context.Background()))
	}()

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}
