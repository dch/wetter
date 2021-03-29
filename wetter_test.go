package wetter_test

import (
	"log"
	"os"
	"testing"
	"wetter"
)

func TestMain(m *testing.M) {
	_, ok := os.LookupEnv("OWM_TOKEN")
	if !ok {
		log.Fatalf("error: you need to provide an OpenWeatherMap API token")
	}
	os.Exit(m.Run())
}

func TestNewClient(t *testing.T) {
	APIKey, _ := os.LookupEnv("OWM_TOKEN")
	client := wetter.NewClient(APIKey)
	if client.APIKey != APIKey {
		t.Error("want client.APIKey to be set")
	}
}

func TestGet(t *testing.T) {
	APIKey, _ := os.LookupEnv("OWM_TOKEN")
	w, err := wetter.GetWeather(APIKey, "Vienna,AT", false)
	if err != nil {
		t.Fatal(err)
	}
	if w == "" {
		t.Error("want non-empty response")
	}
}

func TestTemperatureRanges(t *testing.T) {
	APIKey, _ := os.LookupEnv("OWM_TOKEN")
	c := wetter.NewClient(APIKey)

	conditions, err := c.GetOwmWeather("Vienna,AT")
	if err != nil {
		t.Errorf("wanted error-free response, got %v", err)
	}
	if conditions.TemperatureKelvin < 250 || conditions.TemperatureKelvin > 320 {
		t.Errorf("wanted temperature between 250K and 320K (approx -20°C to +50°C), got %v", conditions)
	}
}
