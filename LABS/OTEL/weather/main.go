package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/OTEL/weather/infra/telemetry"
	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/OTEL/weather/internal/application/command"
	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/OTEL/weather/internal/domain/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	os.Setenv("APP_NAME", "weather-api")
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	c := command.NewCheckZipCommand(service.NewZipCodeService())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdown, err := telemetry.SetupProvider(ctx, os.Getenv("APP_NAME"))
	if err != nil {
		return
	}

	mux := http.NewServeMux()
	t := otel.Tracer("weather")

	// Register handlers.
	mux.HandleFunc("GET /{zipcode}", func(w http.ResponseWriter, r *http.Request) {
		carrier := propagation.HeaderCarrier(r.Header)
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		ctx, span := t.Start(ctx, "weather")
		defer span.End()

		zipcode := r.PathValue("zipcode")

		out, err := c.Handle(ctx, command.CheckZipWeatherCommand{
			ZipCode: zipcode,
		})
		if err != nil {
			if err.Error() == "can not find zipcode" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(out)
	})

	srv := &http.Server{
		Addr:         ":8001",
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
