# Meteocat API Client

[![Go Report Card](https://goreportcard.com/badge/github.com/oscaromeu/meteocat)](https://goreportcard.com/report/github.com/oscaromeu/meteocat)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/itchyny/meteocat/blob/master/LICENSE)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/oscaromeu/meteocat)](https://pkg.go.dev/github.com/oscaromeu/meteocat)

### API Rest client implented in Go to send requests and inspect responses of the Meteocat Rest API

This package allows you to use the API provided by Meteocat to retrieve weather data from its weather stations networks.
The package is still in development and does include support for all the API operations. See the TODO section for
more information.

## What kind of data can I get with Meteocat Go Library ?

Access to forecasts, real-time and historical data from the Meteorological Service of Catalonia

## Get started

### API key

As Meteocat APIs need a valid API key to allow responses, this library won't work if you don't provide one. This stands
for both free and paid (pro) subscription plans. You can signup for a free API key on the Meteocat [website](https://apidocs.meteocat.gencat.cat/section/informacio-general/plans-i-registre/). Please notice that
both subscriptions plan are subject to requests throttling.
### Installation

#### Build from source

`go get github.com/oscaromeu/meteocat`

### Examples

#### Get value of Minimum subsoil temperature at 5 cm at the Viladecans station

```go
package main

import (
	"fmt"
	"github.com/oscaromeu/meteocat"
	"log"
	"os"
)

func main() {

	// execute export METEOCAT_API_KEY=<API_KEY_VALUE> on a shell first
	d, err := meteocat.NewMesurades(os.Getenv("METEOCAT_API_KEY"))
	if err != nil {
		log.Fatalln(err)
	}

	if meteocat.CheckAPIKeyExists(d.Key) == false {
		fmt.Println("ApiKey is not set. ")
	}

	data := meteocat.Data{
		Any: "2021",
		Mes: "01",
		Dia: "06",
	}
	params, _ := meteocat.NewParameters(
		meteocat.OptionCodiEstacio("UG"),
		meteocat.OptionCodiVariable("5"),
		meteocat.OptionData(data),
	)

	// Call MeasurementByDay Method
	d.MeasurementByDay(params)
	fmt.Println(d.Measurements)
}
```

## Documentation

Documentation of the API can be found
at [https://apidocs.meteocat.gencat.cat/documentacio/](https://apidocs.meteocat.gencat.cat/documentacio/).

## TODO
- [ ] Add trace request latency (DNSLookup, TCP Connection and so on)
- [ ] Add example to ingest data on Elasticsearch and visualize it on Kibana
- [ ] Add example to insert data on Influxdb and visualize on Grafana Dashboard
- [ ] Add CLI features 
- [ ] Add support for the following API operations:
    - [x] Mesurades
    - [ ] Predicció
    - [ ] Representatives
    - [ ] Estacions
    - [ ] Estadistics
    - [ ] Càlcul multivariable
    - [ ] Xarxa de Detecció de Descàrregues Elèctriques
    - [ ] Referència
    - [ ] Quotes
- [ ] Full code coverage

## Bug Tracker

Report bug at [Issues・oscaromeu/meteocat - GitHub](https://github.com/oscaromeu/meteocat/issues).

## Author

oscaromeu (https://github.com/oscaromeu)

## License

This software is released under the MIT License, see LICENSE.
