// public interface for wetter API
package main

import (
	"fmt"
	"os"
	"wetter"
)

const owmLocation = "vienna,at"

func main() {
	// retrieve API token, bail if not present
	owmApiToken, ok := os.LookupEnv("OWM_TOKEN")
	if !ok {
		fmt.Fprintf(os.Stderr, "error: you need to provide an OpenWeatherMap API token\n")
		os.Exit(1)
	}
	fmt.Println("Wetter API token:", owmApiToken)

	location := owmLocation

	if len(os.Args) == 2 {
		location = os.Args[1]
	}

	conditions, err := wetter.Weather(owmApiToken, location)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: OpenWeatherMap API response was unparseable: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+#v\n", conditions)
}
