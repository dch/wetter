package wetter_test

import (
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


