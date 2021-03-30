package wetter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)


type client struct {
	APIKey string
	BaseURI string
	HTTPClient *http.Client
}

func NewClient(APIKey string) *client {
	return &client{
		APIKey: APIKey,
		BaseURI: "https://api.openweathermap.org",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *client) GetWeather(location string) (Conditions, error) {
	owmApiURI := fmt.Sprintf("%s/data/2.5/weather?appid=%s&q=%s",
		c.BaseURI,
		c.APIKey,
		url.QueryEscape(location))

	resp, err := c.HTTPClient.Get(owmApiURI)
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
	conditions, err := c.GetWeather(OWMLocation)
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

