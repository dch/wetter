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
	c := wetter.NewClient(APIKey)
	w, err := c.Weather("vienna,AT")
	if err != nil {
		t.Fatal(err)
	}
	if w.Summary == "" {
		t.Error("want non-empty summary")
	}
	if w.TemperatureCentigrade <= -273.15 {
		t.Errorf("want temperature greater than absolute zero")
	}
}