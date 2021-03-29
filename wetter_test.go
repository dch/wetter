package wetter_test

import (
	"os"
	"testing"
	"wetter"
)

func TestNewClient(t *testing.T) {
	APIKey, ok := os.LookupEnv("OWM_TOKEN")
	if !ok {
		t.Error("error: you need to provide an OpenWeatherMap API token")
	}
	client := wetter.NewClient(APIKey)
	if client.APIKey != APIKey {
		t.Error("want client.APIKey to be set")
	}
}

func TestGet(t *testing.T) {
	APIKey, ok := os.LookupEnv("OWM_TOKEN")
	if !ok {
		t.Error("error: you need to provide an OpenWeatherMap API token")
	}
	w, err := wetter.GetWeather(APIKey, "Vienna,AT", false)
	if err != nil {
		t.Fatal(err)
	}
	if w == "" {
		t.Error("want non-empty response")
	}
}
