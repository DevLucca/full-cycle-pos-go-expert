package service

import (
	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/cep-api/internal/domain/entity"
	"github.com/lccmrx/full-cycle-pos-go-expert/LABS/cep-api/internal/domain/vo"
)

type ZipCode struct{}

func NewZipCodeService() *ZipCode {
	return &ZipCode{}
}

func (s *ZipCode) FillTemperaturesFrom(e *entity.ZipCode, unit vo.MeasurementUnit, value float64) {
	temp := vo.NewTemperature(unit, value)
	e.AddTemperature(*temp)

	for unit := range vo.MeasurementUnitsConversions {
		e.AddTemperature(*temp.ToTemperature(unit))
	}
}
