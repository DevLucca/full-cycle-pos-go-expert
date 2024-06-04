package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/OTEL/weather/internal/domain/entity"
	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/OTEL/weather/internal/domain/service"
	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/OTEL/weather/internal/domain/vo"
	"go.opentelemetry.io/otel"
)

type CheckZipCommand struct {
	client http.Client
	svc    *service.ZipCode
}

func NewCheckZipCommand(svc *service.ZipCode) *CheckZipCommand {
	return &CheckZipCommand{
		svc:    svc,
		client: http.Client{},
	}
}

type CheckZipWeatherCommand struct {
	ZipCode string `json:"zip"`
}

type CheckZipWeatherCommandOut struct {
	Zip        string  `json:"zip"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func (c *CheckZipCommand) Handle(ctx context.Context, cmd CheckZipWeatherCommand) (CheckZipWeatherCommandOut, error) {
	ctx, commandSpan := otel.GetTracerProvider().Tracer("weather").Start(ctx, "weather-command")
	defer commandSpan.End()

	var data any

	ctx, zipcodeQuerySpan := otel.GetTracerProvider().Tracer("weather").Start(ctx, "weather-zipcode-query")
	r, _ := http.NewRequest("GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json", cmd.ZipCode), nil)
	resp, err := c.client.Do(r)
	zipcodeQuerySpan.End()
	if err != nil {
		return CheckZipWeatherCommandOut{}, err
	}

	json.NewDecoder(resp.Body).Decode(&data)
	if _, ok := data.(map[string]any)["erro"]; ok {
		return CheckZipWeatherCommandOut{}, errors.New("can not find zipcode")
	}

	escapedCity := url.QueryEscape(data.(map[string]any)["localidade"].(string))

	ctx, weatherQuerySpan := otel.GetTracerProvider().Tracer("weather").Start(ctx, "weather-query")
	r, _ = http.NewRequest("GET", fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=b95773e12eee4a5495943142240505&q=%s", escapedCity), nil)
	resp, err = c.client.Do(r)
	weatherQuerySpan.End()
	if err != nil {
		return CheckZipWeatherCommandOut{}, err
	}

	if resp.StatusCode >= 400 {
		return CheckZipWeatherCommandOut{}, err
	}

	json.NewDecoder(resp.Body).Decode(&data)
	celsiusTemp := data.(map[string]any)["current"].(map[string]any)["temp_c"].(float64)

	e := entity.NewZipCode(entity.ZipCodeParams{
		Id: cmd.ZipCode,
	})

	_, temperatureConversion := otel.GetTracerProvider().Tracer("weather").Start(ctx, "weather-temp-conversion")
	c.svc.FillTemperaturesFrom(e, vo.Celsius, celsiusTemp)
	temperatureConversion.End()

	return CheckZipWeatherCommandOut{
		Zip:        cmd.ZipCode,
		Celsius:    e.GetTemperature(vo.Celsius),
		Fahrenheit: e.GetTemperature(vo.Fahrenheit),
		Kelvin:     e.GetTemperature(vo.Kelvin),
	}, nil
}
