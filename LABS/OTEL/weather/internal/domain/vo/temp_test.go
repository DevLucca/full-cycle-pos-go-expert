package vo

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewTemperature(t *testing.T) {
	type args struct {
		unit  MeasurementUnit
		value float64
	}
	tests := []struct {
		name string
		args args
		want *Temperature
	}{
		{name: "celsius", args: args{unit: Celsius, value: 30.0}, want: &Temperature{MeasurementUnit: Celsius, Value: 30.0}},
		{name: "fahrenheit", args: args{unit: Fahrenheit, value: 30.0}, want: &Temperature{MeasurementUnit: Fahrenheit, Value: 30.0}},
		{name: "kelvin", args: args{unit: Kelvin, value: 30.0}, want: &Temperature{MeasurementUnit: Kelvin, Value: 30.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTemperature(tt.args.unit, tt.args.value)

			if got := got; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTemperature() = %v, want %v", got, tt.want)
			}

			fmt.Println(got.ToTemperature(Fahrenheit))
		})
	}
}

func TestTemperature_ToTemperature(t *testing.T) {
	type fields struct {
		Value           float64
		MeasurementUnit MeasurementUnit
	}
	type args struct {
		unit MeasurementUnit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Temperature
	}{
		{name: "celsius to fahrenheit", fields: fields{Value: 30, MeasurementUnit: Celsius}, args: args{unit: Fahrenheit}, want: &Temperature{MeasurementUnit: Fahrenheit, Value: 86}},
		{name: "celsius to kelvin", fields: fields{Value: 30, MeasurementUnit: Celsius}, args: args{unit: Kelvin}, want: &Temperature{MeasurementUnit: Kelvin, Value: 303.15}},
		{name: "fahrenheit to celsius", fields: fields{Value: 86, MeasurementUnit: Fahrenheit}, args: args{unit: Celsius}, want: &Temperature{MeasurementUnit: Celsius, Value: 30}},
		{name: "fahrenheit to kelvin", fields: fields{Value: 86, MeasurementUnit: Fahrenheit}, args: args{unit: Kelvin}, want: &Temperature{MeasurementUnit: Kelvin, Value: 303.15}},
		{name: "kelvin to celsius", fields: fields{Value: 303.15, MeasurementUnit: Kelvin}, args: args{unit: Celsius}, want: &Temperature{MeasurementUnit: Celsius, Value: 30}},
		{name: "kelvin to fahrenheit", fields: fields{Value: 303.15, MeasurementUnit: Kelvin}, args: args{unit: Fahrenheit}, want: &Temperature{MeasurementUnit: Fahrenheit, Value: 86}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Temperature{
				Value:           tt.fields.Value,
				MeasurementUnit: tt.fields.MeasurementUnit,
			}
			if got := tr.ToTemperature(tt.args.unit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Temperature.ToTemperature() = %v, want %v", got, tt.want)
			}
		})
	}
}
