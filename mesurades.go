package meteocat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// MetadadesVariable is an aggregation of fields to hold the metadata asociated with
// the variables of the Network of Automatic Meteorological Stations (XEMA), integrated into the Network of Meteorological
// Equipment of the Generalitat de Catalunya (Xemec), of the Meteorological Service of Catalonia. Each variable is identified by a code.
type MetadadesVariable struct {
	Codi     int    `json:"codi"`     // Identifier code of the variable
	Nom      string `json:"nom"`      // Identifier name of the variable
	Unitats  string `json:"unitats"`  // Unit of measurement of variables e.g T(ºC),...
	Acronim  string `json:"acronim"`  // Acronym of the variable
	Tipus    string `json:"tipus"`    // Type of variable
	Decimals int    `json:"decimals"` // Number of decimal numbers
}

// MetadadesVariables is a slice which holds the metadata of all variables
type MetadadesVariables []struct{ MetadadesVariable }

// Estats struct holds information of station code when initiated and or finalized.
type Estats []struct {
	Codi      int         `json:"codi"`
	DataInici string      `json:"dataInici"`
	DataFi    interface{} `json:"dataFi"`
}

// BasesTemporals TODO
type BasesTemporals []struct {
	Codi      string `json:"codi"`
	DataInici string `json:"dataInici"`
	DataFi    string `json:"dataFi"`
}

// MetadadesVariablesEstacio is a slice which holds the
// variables metadata of all the data registered bu a station
type MetadadesVariablesEstacio []struct{ MetadadesVariableEstacio }

// MetadadesVariableEstacio is an agreggate type to hold the metadata of the variable data registered in a particular station
type MetadadesVariableEstacio struct {
	MetadadesVariable
	Estats
	BasesTemporals
}

// Lectura is an aggregate type which represents the data registered in the station. This value with a code represents
// a variable, e.g {"codi":5,"lectures":[{"data":"2021-01-06T10:00Z","dataExtrem":"2021-01-06T10:24Z","valor":8.7,"estat":" ","baseHoraria":"SH"}]}
type Lectura struct {
	Data        string `json:"data"`
	Valor       float64 `json:"valor"`
	Estat       string  `json:"estat"`
	BaseHoraria string  `json:"baseHoraria"`
}

// Variable is an agreggate type which represents the variable data registered in a station.
type Variable struct {
	Codi     int       `json:"codi"`
	Lectures []Lectura `json:"lectures"`
}

// Measurements holds the measurements done in a station
type Measurements []struct {
	Codi      string     `json:"codi"`
	Variables []Variable `json:"variables"`
}

// Mesurades holds all the data representations to unmarshall the API responses
type Mesurades struct {
	Variable
	Measurements
	MetadadesVariablesEstacio
	MetadadesVariableEstacio
	MetadadesVariables
	Key          string
	CodiEstacio  string // ?
	CodiVariable string // ?
	*Settings
}

// NewMesurades returns a new MesuradesData pointer with the supplied parameters
func NewMesurades(key string) (*Mesurades, error) {
	c := &Mesurades{
		Settings: NewSettings(),
	}

	c.Key, _ = setKey(key)

	return c, nil
}

// Returns information of a variable for all stations for a given day, if the station code is reported,
// returns the data of the variable for the requested station.
// The API resource is /variables/mesurades/{codi_variable}/{any}/{mes}/{dia}?codiEstacio={codi_estacio} where the
// parameters `codi_variable`, `any`, `mes`, `dia` are mandatory and `codi_estacio` is optional.
// Request example: https://api.meteo.cat/xema/v1/variables/mesurades/32/2017/03/27?codiEstacio=UG
func (m *Mesurades) MeasurementByDay(p *Parameters) error {

	if ValidCodiVariable(p.codiVariable) {
		m.CodiVariable = p.codiVariable
	} else {
		return errVariableUnavailable
	}

	if p.codiEstacio != "" {
		p.codiEstacio = strings.ToUpper(p.codiEstacio)

		if ValidCodiEstacio(p.codiEstacio) {
			m.CodiEstacio = p.codiEstacio
			sFlag = true
			url = fmt.Sprintf(baseURL, fmt.Sprintf("/variables/mesurades/%s/%s/%s/%s?codiEstacio=%s", p.codiVariable, p.Any, p.Mes, p.Dia, p.codiEstacio))
		} else {
			return errEstacioUnavailable
		}

	} else {
		url = fmt.Sprintf(baseURL, fmt.Sprintf("/variables/mesurades/%s/%s/%s/%s", p.codiVariable, p.Any, p.Mes, p.Dia))
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", m.Key)

	resp, err := m.client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if sFlag == true {
		if err = json.NewDecoder(resp.Body).Decode(&m.Variable); err != nil {
			return err
		}
	} else {
		if err = json.NewDecoder(resp.Body).Decode(&m.Measurements); err != nil {
			return err
		}
	}
	// print status line and headers
	return nil
}

// Returns information of all variables for a station for a given day.
// The API resource is /estacions/mesurades/{codiEstacio}/{any}/{mes}/{dia} where the parameters `codiEstacio`  `any`
// `mes` and `dia` are all mandatory. Request example: https://api.meteo.cat/xema/v1/estacions/mesurades/CC/2020/06/16

func (m *Mesurades) MeasurementAllByStation(p *Parameters) error {
	if ValidCodiEstacio(p.codiEstacio) {
		m.CodiEstacio = strings.ToUpper(p.codiEstacio)
	} else {
		return errEstacioUnavailable
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(baseURL, fmt.Sprintf("/estacions/mesurades/%s/%s/%s/%s", strings.ToUpper(p.codiEstacio), p.Any, p.Mes, p.Dia)), nil)

	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", m.Key)

	resp, err := m.client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&m.Measurements); err != nil {
		return err
	}

	return nil
}

// Returns the last measurement of the last 4 hours for all stations of a variable, filtered by station if indicated
// The API resource is /variables/mesurades/{codi_variable}/ultimes?codiEstacio={codi_estacio} where `codi_variable` is
// mandatory and `codi_estacio` is optional. Request example: https://api.meteo.cat/xema/v1/variables/mesurades/5/ultimes?codiEstacio=UG
func (m *Mesurades) MeasurementLast(p *Parameters) error {

	if ValidCodiVariable(p.codiVariable) {
		m.CodiVariable = p.codiVariable
	} else {
		return errVariableUnavailable
	}

	if p.codiEstacio != "" {
		p.codiEstacio = strings.ToUpper(p.codiEstacio)
		if ValidCodiEstacio(p.codiEstacio) {
			m.CodiEstacio = p.codiEstacio
			url = fmt.Sprintf(baseURL, fmt.Sprintf("/variables/mesurades/%s/ultimes?codiEstacio=%s", p.codiVariable, p.codiEstacio))
			sFlag = true
		} else {
			return errEstacioUnavailable
		}

	} else {
		url = fmt.Sprintf(baseURL, fmt.Sprintf("/variables/mesurades/%s/ultimes", p.codiVariable))
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}


	req.Header.Add("X-Api-Key", m.Key)

	resp, err := m.client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if sFlag == true {
		if err = json.NewDecoder(resp.Body).Decode(&m.Variable); err != nil {
			return err
		}
	} else {
		if err = json.NewDecoder(resp.Body).Decode(&m.Measurements); err != nil {
			return err
		}
	}

	return nil
}

// Returns metadata of all variables measured by the station with code specified in the URL, filtered by status and date if specified
// The API resource is /estacions/{codiEstacio}/variables/mesurades/metadades?estat={estat}&data={data} where `estat` and `date` are optional.
// The `estat` parameter describes the state of the station and it can have one of the following values [ope, des, bte]. These values means
// "Operativa", "Baixa temporal" and "Desmantellada" respectively. See https://apidocs.meteocat.gencat.cat/documentacio/dades-de-la-xema/ for more information.
// Note that the 'data' and 'estat' parameters are required together in order to filter the metadata
// Request example https://api.meteo.cat/xema/v1/estacions/UG/variables/mesurades/metadades?estat=ope&data=2017-03-27Z
func (m *Mesurades) MeasurementMetadataAllByStation(p *Parameters) error {

	if ValidCodiEstacio(p.codiEstacio) {
		m.CodiEstacio = p.codiEstacio
	} else {
		return errVariableUnavailable
	}

	dataOk := ValidData(p.Data)
	if p.codiEstat != "" && dataOk == true {
		p.codiEstat = strings.ToLower(p.codiEstat)

		if ValidCodiEstat(p.codiEstat) {
			sFlag = true
			url = fmt.Sprintf(baseURL, fmt.Sprintf("/estacions/%s/variables/mesurades/metadades?estat=%s&data=%s-%s-%sZ", p.codiEstacio, p.codiEstat, p.Any, p.Mes, p.Dia))
		} else {
			return errEstacioUnavailable
		}

	} else {
		url = fmt.Sprintf(baseURL, fmt.Sprintf("/estacions/%s/variables/mesurades/metadades", p.codiEstacio))
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", m.Key)

	resp, err := m.client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if sFlag == true {
		if err = json.NewDecoder(resp.Body).Decode(&m.MetadadesVariableEstacio); err != nil {
			return err
		}
	} else {
		if err = json.NewDecoder(resp.Body).Decode(&m.MetadadesVariablesEstacio); err != nil {
			return err
		}
	}
	// print status line and headers
	return nil
}

// MeasurementMetadataByStation Returns the metadata of the variable with the code specified in the URL that measures the station with the code indicated in the URL
// The API resource is /estacions/{codiEstacio}/variables/mesurades/{codiVariable}/metadades where the parameters 'codiEstacio' and 'codiVariable' are mandatory
// Request example https://api.meteo.cat/xema/v1/estacions/UG/variables/mesurades/3/metadades
func (m *Mesurades) MeasurementMetadataByStation(p *Parameters) error {

	if ValidCodiVariable(p.codiVariable) {
		m.CodiVariable = p.codiVariable
	} else {
		return errVariableUnavailable
	}

	if ValidCodiEstacio(p.codiEstacio) {
		m.CodiEstacio = p.codiEstacio
	} else {
		return errEstacioUnavailable
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(baseURL, fmt.Sprintf("/estacions/%s/variables/mesurades/%s/metadades", p.codiEstacio, p.codiVariable)), nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", m.Key)

	resp, err := m.client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&m.MetadadesVariableEstacio); err != nil {
		return err
	}

	return nil
}

// Returns the metadata of all variables regardless of the stations at which they are measured. The API resource is /variables/mesurades/metadades and
// there are no parameters. Request example https://api.meteo.cat/xema/v1/variables/mesurades/metadades
func (m *Mesurades) MeasurementMetadataAll() error {

	req, err := http.NewRequest("GET", fmt.Sprintf(baseURL, fmt.Sprintf("/variables/mesurades/metadades")), nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", m.Key)

	resp, err := m.client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&m.MetadadesVariables); err != nil {
		return err
	}

	return nil
}

// Returns the metadata of the variable with code indicated in the URL, regardless of the stations in which they are measured.
// The API resource is /variables/mesurades/{codi_variable}/metadades where the parameter 'codi_variable' is mandatory.
// Request example: https://api.meteo.cat/xema/v1/variables/mesurades/1/metadades
func (m *Mesurades) MeasurementMetadataUnique(p *Parameters) error {

	if ValidCodiVariable(p.codiVariable) {
		m.CodiVariable = p.codiVariable
	} else {
		return errVariableUnavailable
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(baseURL, fmt.Sprintf("/variables/mesurades/%s/metadades", p.codiVariable)), nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", m.Key)
	resp, err := m.client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&m.MetadadesVariable); err != nil {
		return err
	}

	return nil
}
