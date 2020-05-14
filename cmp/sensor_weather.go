package cmp

import "github.com/tsybulin/gosha/evt"

// WeatherInfo ...

type WeatherInfo struct {
	Weather []struct {
		Main string `json:"main"`
		Icon string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float32 `json:"temp"`
		FeelsLike float32 `json:"feels_like"`
		TempMin   float32 `json:"temp_min"`
		TempMax   float32 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float32 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Sys struct {
		Sunrise int `json:"sunrise"`
		Sunset  int `json:"sunset"`
	} `json:"sys"`
}

// WeatherSensor ...
type WeatherSensor interface {
	Sensor
	SetAttributes(WeatherInfo)
	SetValues(event evt.Message)
}
