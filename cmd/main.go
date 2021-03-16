// public interface for wetter API
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"wetter"
)

var owmApiBaseURI = "https://api.openweathermap.org"
var owmLocation = "vienna,at"

type OwmResponse struct {
	Weather []struct {
		Main string
	}
	Main struct {
		Kelvin float64 `json:"temp"`
	}
}

func main() {
	// retrieve API token, bail if not present
	owmApiToken, ok := os.LookupEnv("OWM_TOKEN")
	if !ok {
		fmt.Fprintf(os.Stderr, "error: you need to provide an OpenWeatherMap API token\n")
		os.Exit(1)
	}
	fmt.Println("Wetter API token:", owmApiToken)

	owmApiURI := fmt.Sprintf("%s/data/2.5/weather?appid=%s&q=%s",
		owmApiBaseURI,
		owmApiToken,
		url.QueryEscape(owmLocation))

	resp, err := http.Get(owmApiURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: OpenWeatherMap API call returned an error: %v\n", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "error: OpenWeatherMap API response status was unexpected: %s\n", resp.Status)
		os.Exit(1)
	}

	// got 200, so query was valid, and data was available, this is
	// probably parseable JSON now
	// John's trick - read entire body into [] via ioutil.ReadAll
	// look up io.MultiWriter for more details
	var owmResponse OwmResponse
	err = json.NewDecoder(resp.Body).Decode(&owmResponse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: OpenWeatherMap API response was unparseable: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", owmResponse)
}
