package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	client http.Client
	db     *sql.DB
)

func init() {
	client = http.Client{
		Timeout: 200 * time.Millisecond,
	}
	var err error
	db, err = sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/cotacao", Handler)
	slog.Info("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

type response struct {
	Err  string `json:"error,omitempty"`
	Data any    `json:"data,omitempty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jsonEnc := json.NewEncoder(w)

	p, err := getDollarPrice(ctx)
	if err != nil {
		slog.Error("error getting dollar price from API", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		jsonEnc.Encode(response{Err: "error getting dollar price from API"})
		return
	}

	Save(ctx, p)
	jsonEnc.Encode(response{Data: p})
}

type Price struct {
	Bid string `json:"bid"`
}

func getDollarPrice(ctx context.Context) (p Price, err error) {

	res, err := client.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	in := make(map[string]any)
	json.Unmarshal(body, &in)

	body, err = json.Marshal(in["USDBRL"])
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		return
	}

	return
}

func Save(ctx context.Context, p Price) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	res, err := db.ExecContext(ctx, "INSERT INTO price (bid) VALUES (?)", p.Bid)
	if err != nil {
		slog.Error("error saving price", slog.String("error", err.Error()))
		return
	}
	id, _ := res.LastInsertId()
	slog.Info("price saved", slog.String("id", fmt.Sprintf("%v", id)))
}
