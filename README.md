# Wetter

A BSD-licensed command-line interface to the [OpenWeatherMap API] for fun.

It will guess your location from whatever external (internet-facing) IP
address we have, to determine the most appropriate geographic region and
location for guessing your weather.

It supports imperial and metric units, defaulting to metric units, and
you can optionally override it with a suitable location, if the guessed
location is not to your preference, or you want to obtain the weather
elsewhere.

## License

[BSD-2 Clause](https://skunkwerks.github.io/bsd-2-clause-license)

## Pre-requisites

- You'll need to sign up for an [OpenWeatherMap API] key stored in
  the `OWM_TOKEN` environment variable
- internet connectivity for GeoIP lookup and API access

## Usage

```
$ env OWM_TOKEN=... wetter [-i] \
    [-l <integer | city,ISO3166 CountryCode>]


Your nearest weather is Vienna,AT is:

overcast with clouds

4.9°С temperature from 3.9 to 6.1 °С, wind 2.02 m/s. clouds 98 %, 1017 hpa

Geo coords [48.2085, 16.3721]

$ wetter -l 2761369
$ wetter -i -l vienna,at
```
