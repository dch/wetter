// public interface for wetter API
package main

import (
	"fmt"
	"os"
)

func main() {
	// retrieve API token, bail if not present
	// what is the right way to bail this early?
	// can|should we test this?
	owmApi, ok := os.LookupEnv("OWM_TOKEN")
	if !ok {
		fmt.Fprintf(os.Stderr, "error: you need to provide an OpenWeatherMap API token\n")
		os.Exit(1)
	}
	fmt.Println("Wetter API token:", owmApi)
}
