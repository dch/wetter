package wetter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const owmApiBaseURI = "https://api.openweathermap.org"

type Conditions struct {
	Summary           string
	TemperatureKelvin float64
}

type Client struct {
	APIKey string
}

func NewClient(APIKey string) *Client {
	return &Client{
		APIKey: APIKey,
	}
}

type owmResponse struct {
	Weather []struct {
		Main string
	}
	Main struct {
		Kelvin float64 `json:"temp"`
	}
}

func (c Conditions) String(useFahrenheit bool) string {

	celsius := c.TemperatureKelvin - 273.15
	temperature := ""

	if useFahrenheit {
		temperature = fmt.Sprintf("%0.2f °F", celsius*9/5+32)
	} else {
		temperature = fmt.Sprintf("%0.2f °C", celsius)
	}

	return c.Summary + ", " + temperature
}
func Weather(owmApiToken string, owmLocation string) (Conditions, error) {
	c := NewClient(owmApiToken)
	
	conditions, err := c.getOwmWeather(owmLocation)
	if err != nil {
		return Conditions{}, err
	}
	return conditions, nil
}

func (c *Client) getOwmWeather(owmLocation string) (Conditions, error) {
	owmApiURI := fmt.Sprintf("%s/data/2.5/weather?appid=%s&q=%s",
		owmApiBaseURI,
		c.APIKey,
		url.QueryEscape(owmLocation))

	resp, err := http.Get(owmApiURI)
	if err != nil {
		return Conditions{}, fmt.Errorf("error: OpenWeatherMap API call returned an error: %v\n", err)
	}
	defer resp.Body.Close()

	// TLS:(*tls.ConnectionState)(0xc00051d080)}
	// &http.Response{Status:"200 OK", StatusCode:200, Proto:"HTTP/1.1",
	// ProtoMajor:1, ProtoMinor:1,
	// Header:http.Header{"Access-Control-Allow-Credentials":[]string{"true"},
	// "Access-Control-Allow-Methods":[]string{"GET, POST"},
	// "Access-Control-Allow-Origin":[]string{"*"},
	// "Connection":[]string{"keep-alive"}, "Content-Length":[]string{"481"},
	// "Content-Type":[]string{"application/json; charset=utf-8"},
	// "Date":[]string{"Tue, 02 Feb 2021 12:45:02 GMT"},
	// "Server":[]string{"openresty"},
	// "X-Cache-Key":[]string{"/data/2.5/weather?q=wellington"}},
	// Body:(*http.bodyEOFSignal)(0xc000528080), ContentLength:481,
	// TransferEncoding:[]string(nil), Close:false, Uncompressed:false,
	// Trailer:http.Header(nil), Request:(*http.Request)(0xc000144000),

	if resp.StatusCode != http.StatusOK {
		return Conditions{}, fmt.Errorf("error: OpenWeatherMap API response status was unexpected: %s\n", resp.Status)
	}

	// got 200, so query was valid, and data was available, this is
	// probably parseable JSON now
	// John's trick - read entire body into [] via ioutil.ReadAll
	// look up io.MultiWriter for more details
	var response owmResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return Conditions{}, fmt.Errorf("OpenWeatherMap API response was unparseable: %v\n", err)
	}
		if len(response.Weather) < 1 {
		return Conditions{}, fmt.Errorf("weather array was empty: %+v", response)
	}

	return Conditions{
		Summary:               response.Weather[0].Main,
		TemperatureCentigrade: response.Main.Kelvin - 273.15,
	}, nil
}