package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	client http.Client
)

type response struct {
	Err  string `json:"error,omitempty"`
	Data any    `json:"data,omitempty"`
}

type Price struct {
	Bid string `json:"bid"`
}

func init() {
	client = http.Client{
		Timeout: 300 * time.Millisecond,
	}
}

func main() {
	res, err := client.Get("http://localhost:8080/cotacao")
	if err != nil {
		slog.Error("error getting dollar price from API", slog.String("error", err.Error()))
		return
	}

	if res.StatusCode != http.StatusOK {
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var r response
	json.Unmarshal(body, &r)
	if r.Err != "" {
		return
	}

	os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dolar: %s\n", r.Data.(map[string]any)["bid"])), 0644)
}
