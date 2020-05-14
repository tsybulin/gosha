package intr

import (
	"fmt"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type weather struct {
	cmp.Component
	value             interface{}
	icon              string
	deviceClass       string
	unitOfMeasurement string
	temp              float32
	feelsLike         float32
	tempMin           float32
	tempMax           float32
	pressure          int
	humidity          int
	visibility        int
	windSpeed         float32
	windDeg           int
	sunrise           int
	sunset            int
}

func (w *weather) GetValue() interface{} {
	return w.value
}

func (w *weather) SetValue(v interface{}) {
	w.value = v
	// log.Println("--- value", w.GetID(), w.value)
}

func (w *weather) SetAttributes(wi cmp.WeatherInfo) {
	w.value = wi.Weather[0].Main
	w.icon = wi.Weather[0].Icon
	w.temp = wi.Main.Temp
	w.feelsLike = wi.Main.FeelsLike
	w.tempMin = wi.Main.TempMin
	w.tempMax = wi.Main.TempMax
	w.pressure = wi.Main.Pressure
	w.humidity = wi.Main.Humidity
	w.visibility = wi.Visibility
	w.windSpeed = wi.Wind.Speed
	w.windDeg = wi.Wind.Deg
	w.sunrise = wi.Sys.Sunrise
	w.sunset = wi.Sys.Sunset
}

func (w *weather) SetValues(event evt.Message) {
	s := event.Event.Data.NewState
	w.value = s.State

	w.icon = s.Attributes["icon"].(string)
	w.temp = s.Attributes["temp"].(float32)
	w.feelsLike = s.Attributes["feels_like"].(float32)
	w.tempMin = s.Attributes["temp_min"].(float32)
	w.tempMax = s.Attributes["temp_max"].(float32)
	w.pressure = s.Attributes["pressure"].(int)
	w.humidity = s.Attributes["humidity"].(int)
	w.visibility = s.Attributes["visibility"].(int)
	w.windSpeed = s.Attributes["wind_speed"].(float32)
	w.windDeg = s.Attributes["wind_deg"].(int)
	w.sunrise = s.Attributes["sunrise"].(int)
	w.sunset = s.Attributes["sunset"].(int)
}

func (w *weather) GetDeviceClass() string {
	return w.deviceClass
}

func (w *weather) GetUnitOfMeasurement() string {
	return w.unitOfMeasurement
}

func (w *weather) GetState() evt.State {
	state := w.Component.GetState()

	state.Attributes["device_class"] = w.deviceClass
	state.State = fmt.Sprint(w.value)
	state.Attributes["unit_of_measurement"] = w.unitOfMeasurement

	state.Attributes["icon"] = w.icon
	state.Attributes["temp"] = w.temp
	state.Attributes["feels_like"] = w.feelsLike
	state.Attributes["temp_min"] = w.tempMin
	state.Attributes["temp_max"] = w.tempMax
	state.Attributes["pressure"] = w.pressure
	state.Attributes["humidity"] = w.humidity
	state.Attributes["visibility"] = w.visibility
	state.Attributes["wind_speed"] = w.windSpeed
	state.Attributes["wind_deg"] = w.windDeg
	state.Attributes["sunrise"] = w.sunrise
	state.Attributes["sunset"] = w.sunset

	return state
}

// NewWeather ...
func NewWeather(cfg map[string]string) cmp.WeatherSensor {
	return &weather{
		Component:         NewComponent(cmp.DomainWeather, cfg["sensor"], "weather"),
		deviceClass:       "weather",
		unitOfMeasurement: "weather",
	}
}
