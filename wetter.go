package wetter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const owmApiBaseURI = "https://api.openweathermap.org"

type OWMResponse struct {
	Weather []struct {
		Main string
	}
	Main struct {
		Kelvin float64 `json:"temp"`
	}
}


type Conditions struct{
	Summary string
	TemperatureCentigrade float64
}

type Client struct{
	APIKey string
}

func NewClient(APIKey string) *Client {
	return &Client{
		APIKey: APIKey,
	}
}

func (c *Client) Weather(location string) (Conditions, error) {
	owmApiURI := fmt.Sprintf("%s/data/2.5/weather?appid=%s&q=%s",
		owmApiBaseURI,
		c.APIKey,
		url.QueryEscape(location))
	resp, err := http.Get(owmApiURI)
	if err != nil {
		return Conditions{}, fmt.Errorf("OpenWeatherMap API call returned an error: %w\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Conditions{}, fmt.Errorf("OpenWeatherMap API response status was unexpected: %s\n", resp.Status)
	}

	// got 200, so query was valid, and data was available, this is
	// probably parseable JSON now
	// John's trick - read entire body into [] via ioutil.ReadAll
	// look up io.MultiWriter for more details
	var owmResponse OWMResponse
	err = json.NewDecoder(resp.Body).Decode(&owmResponse)
	if err != nil {
		return Conditions{}, fmt.Errorf("OpenWeatherMap API response was unparseable: %v\n", err)
	}
	if len(owmResponse.Weather) < 1 {
		return Conditions{}, fmt.Errorf("weather array was empty: %+v", owmResponse)
	}
	return Conditions{
		Summary: owmResponse.Weather[0].Main,
		TemperatureCentigrade: owmResponse.Main.Kelvin - 273.15,
	}, nil
}

