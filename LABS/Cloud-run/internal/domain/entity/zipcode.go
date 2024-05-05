package entity

import "github.com/lccmrx/full-cycle-pos-go-expert/LABS/cep-api/internal/domain/vo"

type ZipCodeParams struct {
	Id string
}

type ZipCode struct {
	Id           string
	Temperatures map[vo.MeasurementUnit]vo.Temperature
}

func NewZipCode(params ZipCodeParams) *ZipCode {
	return &ZipCode{
		Id:           params.Id,
		Temperatures: make(map[vo.MeasurementUnit]vo.Temperature),
	}
}

func (z *ZipCode) AddTemperature(temp vo.Temperature) {
	z.Temperatures[temp.MeasurementUnit] = temp
}

func (z *ZipCode) GetTemperature(unit vo.MeasurementUnit) float64 {
	if value, ok := z.Temperatures[unit]; ok {
		return value.Value
	}

	return 0
}
