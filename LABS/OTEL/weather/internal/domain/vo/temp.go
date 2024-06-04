package vo

import "fmt"

type MeasurementUnit string

const (
	Celsius    MeasurementUnit = "C"
	Fahrenheit MeasurementUnit = "F"
	Kelvin     MeasurementUnit = "K"
)

var MeasurementUnitsConversions = map[MeasurementUnit]map[MeasurementUnit]func(v float64) float64{
	Celsius: {
		Fahrenheit: func(v float64) float64 {
			return v*1.8 + 32
		},
		Kelvin: func(v float64) float64 {
			return v + 273.15
		},
	},
	Fahrenheit: {
		Celsius: func(v float64) float64 {
			return (v - 32) / 1.8
		},
		Kelvin: func(v float64) float64 {
			return (v-32)*5.0/9.0 + 273.15
		},
	},
	Kelvin: {
		Celsius: func(v float64) float64 {
			return v - 273.15
		},
		Fahrenheit: func(v float64) float64 {
			return (v-273.15)*9.0/5.0 + 32
		},
	},
}

type Temperature struct {
	Value           float64
	MeasurementUnit MeasurementUnit
}

func NewTemperature(unit MeasurementUnit, value float64) *Temperature {
	t := Temperature{
		Value:           value,
		MeasurementUnit: unit,
	}

	return &t
}

func (t *Temperature) String() string {
	return fmt.Sprintf("%.2f%sº", t.Value, t.MeasurementUnit)
}

func (t *Temperature) ToTemperature(unit MeasurementUnit) *Temperature {
	if conversion, ok := MeasurementUnitsConversions[t.MeasurementUnit][unit]; ok {
		return NewTemperature(unit, conversion(t.Value))
	}
	return t
}
