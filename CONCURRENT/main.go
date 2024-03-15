package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type FastestAPI struct {
	API      string
	Res      *http.Response
	TimeTook time.Duration
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("CEP não informado")
		os.Exit(1)
		return
	}
	var cep string = os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	starter := make(chan bool)
	fastestCh := make(chan FastestAPI)

	apis := map[string]*http.Request{
		"viacep": func() *http.Request {
			req, _ := http.NewRequest("GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep), nil)
			return req
		}(),
		"brasilapi": func() *http.Request {
			req, _ := http.NewRequest("GET", fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), nil)
			return req
		}(),
	}

	wg := sync.WaitGroup{}

	wg.Add(len(apis))

	for api, req := range apis {
		go func() {
			wg.Done()
			<-starter
			startTime := time.Now()
			fmt.Printf("calling `%s` \t-----\t started at %s\n", api, startTime.Format("15:04:05.000000000"))
			res, err := call(req)
			if err != nil {
				return
			}

			fastestCh <- FastestAPI{
				API:      api,
				Res:      res,
				TimeTook: time.Since(startTime),
			}
		}()
	}

	wg.Wait()
	close(starter)

	select {
	case fastest := <-fastestCh:
		var data any
		json.NewDecoder(fastest.Res.Body).Decode(&data)

		json2, _ := json.MarshalIndent(data, "", "  ")
		fmt.Printf("API mais rápida: %s | Resposta em: %dms\nDados: %s\n", fastest.API, fastest.TimeTook.Milliseconds(), string(json2))
		os.Exit(0)
	case <-ctx.Done():
		fmt.Println("timeout")
		os.Exit(1)
	}
}

func call(req *http.Request) (res *http.Response, err error) {
	res, err = http.DefaultClient.Do(req)

	if err != nil || res.StatusCode != 200 {
		return nil, errors.New("api returned error or status code different than 200")
	}

	return
}
