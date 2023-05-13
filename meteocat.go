package meteocat

import (
	"errors"
	"fmt"
	"net/http"
)

var errEstacioUnavailable = errors.New("station code unavailable")
var errVariableUnavailable = errors.New("variable code unavailable")
var errInvalidKey = errors.New("invalid api key")
var errInvalidOption = errors.New("invalid option")
var errInvalidHttpClient = errors.New("invalid http client")

// DataUnits represents the character chosen to represent the temperature notation
// var DataUnits = map[string]string{"C": "metric"}
var (
	baseURL = "https://api.meteo.cat/xema/v1%s"
	url     string // helper variable to hold url values
	sFlag   bool   // helper variable to check where to handle JSON that varies between an array or a single item
	// Its worth to note that the single item corresponds to a lecture of a single station whereas an array is a lecture
	// containing data of all stations. See the endpoint of MeasurementByDay function.

)

// Config will hold default settings
type Config struct {
	APIKey string // API Key for connecting to the OWM
}

// Data struct holds the time settings in general the time format will be YYYY/MM/D or YYYY-MM-DZ but this is transparent
// for the final user.
type Data struct {
	Any string
	Mes string
	Dia string
}

type TimeDate struct {
	Hour         string
	Minute       string
	Seconds      string
	Milliseconds string
}

// Parameters holds all the options to be passed in to the methods
type Parameters struct {
	codiEstacio  string // should reference a key in the CodisEstacions map
	codiVariable string // should reference a key in the CodisVariables map
	codiEstat    string // should reference a key in the CodisEstat map
	Data
	TimeDate
}

// NewParameters generates a new Parameters config.
func NewParameters(options ...func(*Parameters) error) (*Parameters, error) {
	p := &Parameters{}

	// Default values...
	p.codiEstacio = ""
	p.codiVariable = ""
	p.codiEstat = ""
	p.Any = ""
	p.Mes = ""
	p.Dia = ""
	// Option paremeters values:
	for _, op := range options {
		err := op(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

// OptionCodiEstat is a helper function to set up the value of CodiEstat to be passed in Parameters struct
func OptionCodiEstat(codiEstat string) func(p *Parameters) error {
	return func(p *Parameters) error {
		p.codiEstat = codiEstat
		return nil
	}
}

// OptionCodiEstacio is a helper function to set up the value of CodiEstacio to be passed in Parameters struct
func OptionCodiEstacio(codiEstacio string) func(p *Parameters) error {
	return func(p *Parameters) error {
		p.codiEstacio = codiEstacio
		return nil
	}
}

// OptionCodiVariable is a helper function to set up the value of CodiVariable to be passed in Parameters struct
func OptionCodiVariable(codiVariable string) func(p *Parameters) error {
	return func(p *Parameters) error {
		p.codiVariable = codiVariable
		return nil
	}
}

// OptionData is a helper function to set up the value of Data to be passed in Parameters struct
func OptionData(d Data) func(p *Parameters) error {
	return func(p *Parameters) error {
		p.Any = d.Any
		p.Mes = d.Mes
		p.Dia = d.Dia
		return nil
	}
}

// Option TimeData is a helper function to set up the value of Timedate to be passed in Parameters struct
func OptionTimeDate(d TimeDate) func(p *Parameters) error {
	return func(p *Parameters) error {
		p.Hour = d.Hour
		p.Minute = d.Minute
		p.Seconds = d.Seconds
		p.Milliseconds = d.Milliseconds
		return nil
	}
}

// APIError returned on failed API calls.
type APIError struct {
	Message string `json:"message"`
	COD     string `json:"cod"`
}

// ApiKey setter function to be passed in the Settings struct, necessary to perform the request
func setKey(key string) (string, error) {
	if err := ValidAPIKey(key); err != nil {
		return "", err
	}
	return key, nil
}

// ValidData validates that we set a correct data
func ValidData(d Data) bool {
	if d.Mes != "" && d.Any != "" && d.Dia != "" {
		return true
	} else if d.Mes == "" || d.Any == "" || d.Dia == "" {
		fmt.Println("To use the data parameter all the fields must be filled. Leave empty if its not used")
	}
	return false
}

// ValidAPIKey makes sure that the key given is a valid one
func ValidAPIKey(key string) error {
	if len(key) != 40 {
		return errors.New("invalid key")
	}
	return nil
}

// ValidCodiEstat makes sure the string passed in is an
// acceptable estat code.
func ValidCodiEstat(c string) bool {
	for d := range CodisEstat {
		if c == d {
			return true
		}
	}
	return false
}

// ValidCodiEstacio makes sure the string passed in is an
// acceptable station code.
func ValidCodiEstacio(c string) bool {
	for d := range CodisEstacions {
		if c == d {
			return true
		}
	}
	return false
}

// ValidCodiVariable makes sure the string passed in is an
// acceptable variable code.
func ValidCodiVariable(c string) bool {
	for d := range CodisVariables {
		if c == d {
			return true
		}
	}
	return false
}

// CheckAPIKeyExists will see if an API key has been set.
func CheckAPIKeyExists(apiKey string) bool { return len(apiKey) > 1 }

// Settings holds the client settings
type Settings struct {
	client *http.Client
	req    *http.Request

	//cr *resty.Client
}

// NewSettings returns a new Setting pointer with default http client.
func NewSettings() *Settings {
	return &Settings{
		client: http.DefaultClient,
	}
}

//func OptionURL(url string) Option {
//	return func(s *Settings) error {
//		s.url = url
//		return nil
//	}
//}

// Optional client settings
type Option func(s *Settings) error

// // WithHttpClient sets custom http client when creating a new Client.
func WithHttpClient(c *http.Client) Option {
	return func(s *Settings) error {
		if c == nil {
			return errInvalidHttpClient
		}
		s.client = c
		return nil
	}
}

// setOptions sets Optional client settings to the Settings pointer
func setOptions(settings *Settings, options []Option) error {
	for _, option := range options {
		if option == nil {
			return errInvalidOption
		}
		err := option(settings)
		if err != nil {
			return err
		}
	}
	return nil
}
