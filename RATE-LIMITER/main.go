package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/lccmrx/gorlim"
	"github.com/lccmrx/gorlim/backend"
)

func main() {
	ratelimit := gorlim.New(
		backend.NewRedis(os.Getenv("REDIS_ADDR")),
	)

	http.Handle("GET /", gorlim.Wrap(ratelimit,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var out any
				w.Header().Set("Content-Type", "application/json")
				json.Unmarshal([]byte(`{"ok": true}`), &out)
				json.NewEncoder(w).Encode(out)
			}),
	))

	log.Println("listening on port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
