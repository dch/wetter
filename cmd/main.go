// public interface for wetter API
package main

import (
	"flag"
	"fmt"
	"os"
	"wetter"
)

const defaultLocation = "vienna,at"

// the default is, surprisingly, Celsius, and not Kelvin
var useFahrenheit bool
var owmLocation string

func main() {

	flag.BoolVar(&useFahrenheit, "f", false, "sets display temperature in degrees Fahrenheit (default is Celcius)")
	flag.StringVar(&owmLocation, "l", defaultLocation, "specify a city or location to query temperature")
	flag.Parse()

	if useFahrenheit {
		fmt.Println("We are using Fahrenheit, mein Herr")
	} else {
		fmt.Println("We are using Celsius, my Lord")
	}

	// retrieve API token, bail if not present
	owmApiToken, ok := os.LookupEnv("OWM_TOKEN")
	if !ok {
		fmt.Fprintf(os.Stderr, "error: you need to provide an OpenWeatherMap API token\n")
		os.Exit(1)
	}
	fmt.Println("Wetter API token:", owmApiToken)

	conditions, err := wetter.GetWeather(owmApiToken, owmLocation, useFahrenheit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: OpenWeatherMap API response was inconceivable! %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", owmLocation, conditions)
}
