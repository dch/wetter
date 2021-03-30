//+build integration

package wetter_test

import (
	"os"
	"testing"
	"wetter"
)

func TestGet(t *testing.T) {
	w, err := wetter.GetWeather(getToken(t), "Vienna,AT", false)
	if err != nil {
		t.Fatal(err)
	}
	if w == "" {
		t.Error("want non-empty response")
	}
}

func TestTemperatureRanges(t *testing.T) {
	c := wetter.NewClient(getToken(t))
	conditions, err := c.GetOwmWeather("Vienna,AT")
	if err != nil {
		t.Fatalf("wanted error-free response, got %v", err)
	}
	if conditions.TemperatureKelvin == 0 {
		t.Error("wanted temperature to be set, but was empty")
	}
}

func getToken(t *testing.T) string {
	t.Helper()
	token, ok := os.LookupEnv("OWM_TOKEN")
	if !ok {
		t.Fatalf("error: you need to provide an OpenWeatherMap API token")
	}
	return token
}
