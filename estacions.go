package meteocat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Coordenades struct holds georeference information of the stations
type Coordenades struct {
	Latitud  float64 `json:"latitud"`  // Latitude expressed in decimal degrees. WSG84 reference system
	Longitud float64 `json:"Longitud"` // Longitude expressed in decimal degrees. WSG84 reference system
}

// Municipi struct holds information on which municipi is located the station.
type Municipi struct {
	Codi string `json:"codi"` // INE code of the municipi
	Nom  string `json:"nom"`  // Name of the municipi
}

// Comarca struct holds information on which Comarca is located the station.
type Comarca struct {
	Codi int    `json:"codi"` // Comarca identification code
	Nom  string `json:"nom"`  // Name of the comarca
}

// Provincia struct holds information on which Provincia is located the station.
type Provincia struct {
	Codi int    `json:"codi"` // Provincia identification code
	Nom  string `json:"nom"`  // Name of the provincia
}

// Xarxa struct holds information on which Xarxa is located the station. The xarxa field
// refers to the stations of the Network of Automatic Meteorological Stations (XEMA) of Catalonia
type Xarxa struct {
	Codi int    `json:"codi"` // Xarxa identification code
	Nom  string `json:"nom"`  // Name of the Network
}

// Aggregation of fields and structs to unmarshal the responses
type MetadadesEstacions struct {
	Codi        string      `json:"codi"`        // Identification code for each automatic weather station (EMA)
	Nom         string      `json:"nom"`         // Name of the EMA
	Tipus       string      `json:"tipus"`       // Type of station typically automatic
	Coordenades Coordenades `json:"coordenades"` // Georeference data of the station
	Emplacament string      `json:"emplacament"` // Descriptive name of where the station is located e.g Planters Gusi, ctra. antiga de València, km 14
	Altitud     float64     `json:"altitud"`     // Altitude in meters of the station above the sea level.
	Municipi    Municipi    `json:"municipi"`    // Municipi
	Comarca     Comarca     `json:"comarca"`     // Comarca
	Provincia   Provincia   `json:"provincia"`   // Provincia
	Xarxa       Xarxa       `json:"xarxa"`       // Network tipically XEMA
	Estats      Estats      `json:"estats"`
}

// Mesurades holds all the data representations to unmarshall the API responses
type Estacions struct {
	MetadadesEstacions
	Key          string
	CodiEstacio  string // ?
	CodiVariable string // ?
	*Settings
}

// Returns a list of metadata from all stations. If settings are specified, filters by specified status and date
// The API resource is /estacions/metadades?estat={estat}&data={data} where the
// parameters `estat` and data optional.
// Request example: https://api.meteo.cat/xema/v1/estacions/metadades?estat=ope&data=2017-03-27Z
func (e *Estacions) StationsAll(p *Parameters) error {

	dataOk := ValidData(p.Data)

	if p.codiEstat != "" && dataOk == true {
		p.codiEstat = strings.ToLower(p.codiEstat)

		if ValidCodiEstat(p.codiEstat) {
			sFlag = true
			url = fmt.Sprintf(baseURL, fmt.Sprintf("/estacions/metadades?estat=%s&data=%s-%s-%sZ", p.codiEstat, p.Any, p.Mes, p.Dia))
		} else {
			return errEstacioUnavailable
		}

	} else {
		url = fmt.Sprintf(baseURL, fmt.Sprintf("/estacions/metadades"))
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", e.Key)

	resp, err := e.client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	//if sFlag == true {
	if sFlag {
		if err = json.NewDecoder(resp.Body).Decode(&e.MetadadesEstacions); err != nil {
			return err
		}
	} else {
		if err = json.NewDecoder(resp.Body).Decode(&e.MetadadesEstacions); err != nil {
			return err
		}
	}
	// print status line and headers
	return nil
}
