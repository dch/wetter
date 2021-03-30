package wetter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const OWMAPIBaseURI = "https://api.openweathermap.org"

type client struct {
	APIKey string
}

func NewClient(APIKey string) *client {
	return &client{
		APIKey: APIKey,
	}
}

func (c *client) GetOwmWeather(owmLocation string) (Conditions, error) {
	owmApiURI := fmt.Sprintf("%s/data/2.5/weather?appid=%s&q=%s",
		OWMAPIBaseURI,
		c.APIKey,
		url.QueryEscape(owmLocation))

	resp, err := http.Get(owmApiURI)
	if err != nil {
		return Conditions{}, fmt.Errorf("error: OpenWeatherMap API call returned an error: %v\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Conditions{}, fmt.Errorf("OpenWeatherMap API response status was unexpected: %s", resp.Status)
	}

	// got 200, so query was valid, and data was available, this is
	// probably parseable JSON now
	// John's trick - read entire body into [] via ioutil.ReadAll
	// look up io.MultiWriter for more details
	var response owmResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return Conditions{}, fmt.Errorf("OpenWeatherMap API response was unparseable: %v", err)
	}
	if len(response.Weather) < 1 {
		return Conditions{}, fmt.Errorf("weather array was empty: %+v", response)
	}

	return Conditions{
		Summary:           response.Weather[0].Main,
		TemperatureKelvin: response.Main.Kelvin,
	}, nil
}

type Conditions struct {
	Summary           string
	TemperatureKelvin float64
}

func (c Conditions) Format(useFahrenheit bool) string {
	celsius := c.TemperatureKelvin - 273.15
	temperature := ""

	if useFahrenheit {
		temperature = fmt.Sprintf("%0.2f °F", celsius*9/5+32)
	} else {
		temperature = fmt.Sprintf("%0.2f °C", celsius)
	}

	return c.Summary + ", " + temperature
}

func GetWeather(OWMAPIToken string, OWMLocation string, useFahrenheit bool) (string, error) {
	conditions, err := Weather(OWMAPIToken, OWMLocation)
	if err != nil {
		return "", err
	}
	return conditions.Format(useFahrenheit), nil
}

func Weather(OWMAPIToken string, OWMLocation string) (Conditions, error) {
	c := NewClient(OWMAPIToken)
	conditions, err := c.GetOwmWeather(OWMLocation)
	if err != nil {
		return Conditions{}, err
	}
	return conditions, nil
}

type owmResponse struct {
	Weather []struct {
		Main string
	}
	Main struct {
		Kelvin float64 `json:"temp"`
	}
}

