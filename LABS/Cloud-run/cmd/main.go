package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/cep-api/internal/application/command"
	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/cep-api/internal/domain/service"
)

func main() {
	c := command.NewCheckZipCommand(service.NewZipCodeService())

	http.HandleFunc("/{zipcode}", func(w http.ResponseWriter, r *http.Request) {
		zipcode := r.PathValue("zipcode")

		if matched, _ := regexp.Match(`^\d{8}$`, []byte(zipcode)); !matched {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		out, err := c.Handle(r.Context(), command.CheckZipWeatherCommand{
			ZipCode: zipcode,
		})
		if err != nil {
			if err.Error() == "can not find zipcode" {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(out)
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
