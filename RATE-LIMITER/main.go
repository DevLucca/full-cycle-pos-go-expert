package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lccmrx/gorlim"
	"github.com/lccmrx/gorlim/backend"
)

func main() {
	ratelimit := gorlim.New(
		backend.NewRedis("redis:6379"),
		gorlim.WithRequestLimit(5),
		gorlim.WithTimeframe(gorlim.RPM),
	)

	http.Handle("GET /", gorlim.Wrap(ratelimit,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var out any
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
