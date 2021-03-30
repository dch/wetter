package wetter_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"wetter"
)

func TestNewClient(t *testing.T) {
	wantToken := "dummy token"
	client := wetter.NewClient(wantToken)
	if client.APIKey != wantToken {
		t.Errorf("want %q, got %q", wantToken, client.APIKey)
	}
}

func TestGet(t *testing.T) {
	called := false
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		f, err := os.Open("testdata/vienna.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(w, f)
	}))
	client := wetter.NewClient("dummyToken")
	client.BaseURI = ts.URL
	client.HTTPClient = ts.Client()
	w, err := client.GetWeather("Vienna,AT")
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("server not called")
	}
	wantSummary := "Clouds"
	if wantSummary != w.Summary {
		t.Errorf("want %q, got %q", wantSummary, w.Summary)
	}
	wantTemp := 279.13
	if wantTemp != w.TemperatureKelvin {
		t.Errorf("want %.1f, got %.1f", wantTemp, w.TemperatureKelvin)
	}
}
