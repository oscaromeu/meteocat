package meteocat

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"net/http/httptrace"
	"time"
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
	//baseURL="http://localhost:3000%s"
	openDataURL = "https://analisi.transparenciacatalunya.cat/resource/nzvn-apee.json%s"
	url string // helper variable to hold url values
	sFlag bool // helper variable to check where to handle JSON that varies between an array or a single item

)

// CodisEstat holds all the status that a station can have
var CodisEstat = map[string]string{
	"ope": "Operativa",      // Operational
	"des": "Desmantellada",  // Dismantled
	"bte": "Baixa temporal", // Not operational temporary
}

// CodisEstacions holds all stations in ope state to be used
var CodisEstacions = map[string]string{
	"CC": "Orís",
	"CD": "la Seu d'Urgell - Bellestar",
	"CE": "els Hostalets de Pierola",
	"CG": "Molló - Fabert",
	"CI": "Sant Pau de Segúries",
	"CJ": "Organyà",
	"CL": "Sant Salvador de Guardiola",
	"CP": "Sant Romà d'Abella",
	"CQ": "Vilanova de Meià",
	"CR": "la Quar",
	"CT": "el Pont de Suert",
	"CU": "Vielha",
	"CW": "l'Espluga de Francolí",
	"CY": "Muntanyola",
	"C6": "Castellnou de Seana",
	"C7": "Tàrrega",
	"C8": "Cervera",
	"C9": "Mas de Barberans",
	"DB": "el Perelló",
	"DF": "la Bisbal d'Empordà",
	"DG": "Núria (1.971 m)",
	"DI": "Font-rubí",
	"DJ": "Banyoles",
	"DK": "Torredembarra",
	"DL": "Illa de Buda",
	"DN": "Anglès",
	"DO": "Castell d'Aro",
	"DP": "Das - Aeròdrom",
	"DQ": "Vila-rodona",
	"D1": "Margalef",
	"D2": "Vacarisses",
	"D3": "Vallirana",
	"D4": "Roses",
	"D5": "Barcelona - Observatori Fabra",
	"D6": "Portbou",
	"D7": "Vinebre",
	"D8": "Horta de Sant Joan",
	"D9": "el Vendrell",
	"H1": "Òdena",
	"J5": "Pantà de Darnius - Boadella",
	"KE": "Pantà de Sau",
	"KP": "Fogars de la Selva",
	"KX": "la Roca del Vallès - ETAP Cardedeu",
	"MQ": "Cardona",
	"MR": "Pantà de Siurana",
	"MS": "Castellar de n'Hug - el Clot del Moro",
	"MV": "Guixers - Valls",
	"MW": "Navès",
	"M6": "Sant Joan de les Abadesses",
	"UA": "l'Ametlla de Mar",
	"UB": "la Tallada d'Empordà",
	"UC": "Monells",
	"UE": "Torroella de Montgrí",
	"UF": "PN del Garraf - el Rascler",
	"UG": "Viladecans",
	"UH": "el Montmell",
	"UI": "Gisclareny",
	"UJ": "Santa Coloma de Queralt",
	"UK": "Sant Pere de Ribes - PN del Garraf",
	"UM": "la Granadella",
	"UN": "Cassà de la Selva",
	"UO": "Fornells de la Selva",
	"UP": "Cabrils",
	"UQ": "Dosrius - PN Montnegre Corredor",
	"US": "Alcanar",
	"UU": "Amposta",
	"UW": "els Alfacs",
	"UX": "Ulldecona - els Valentins",
	"UY": "Os de Balaguer - el Monestir d'Avellanes",
	"U1": "Cabanes",
	"U2": "Sant Pere Pescador",
	"U3": "Sant Martí Sarroca",
	"U4": "Castellnou de Bages",
	"U6": "Vinyols i els Arcs",
	"U7": "Aldover",
	"U9": "l'Aldea",
	"VA": "Ascó",
	"VB": "Benissanet",
	"VC": "Pantà de Riba-roja",
	"VD": "el Canós",
	"VE": "Aitona",
	"VH": "Gimenells",
	"VK": "Raimat",
	"VM": "Vilanova de Segrià",
	"VN": "Vilobí d'Onyar",
	"VO": "Lladurs",
	"VP": "Pinós",
	"VQ": "Constantí",
	"VS": "Lac Redon (2.247 m)",
	"VU": "Rellinars",
	"VV": "Sant Llorenç Savall",
	"VX": "Tagamanent - PN del Montseny",
	"VY": "Nulles",
	"VZ": "Espolla",
	"V1": "Vallfogona de Balaguer",
	"V3": "Gurb",
	"V4": "Montesquiu",
	"V5": "Perafita",
	"V8": "el Poal",
	"WA": "Oliola",
	"WB": "Albesa",
	"WC": "Golmés",
	"WD": "Batea",
	"WE": "Vilanova del Vallès",
	"WG": "Algerri",
	"WI": "Maials",
	"WJ": "el Masroig",
	"WK": "Alfarràs",
	"WL": "Sant Martí de Riucorb",
	"WM": "Santuari de Queralt",
	"WN": "Montserrat - Sant Dimes",
	"WO": "la Bisbal del Penedès",
	"WP": "Canaletes",
	"WQ": "Montsec d'Ares (1.572 m)",
	"WR": "Torroja del Priorat",
	"WS": "Viladrau",
	"WT": "Malgrat de Mar",
	"WU": "Badalona - Museu",
	"WV": "Guardiola de Berguedà",
	"WW": "Artés",
	"WX": "Camarasa",
	"WY": "Sant Sadurní d'Anoia",
	"WZ": "Cunit",
	"W1": "Castelló d'Empúries",
	"W4": "la Granada",
	"W5": "Oliana",
	"W8": "Blancafort",
	"W9": "la Vall d'en Bas",
	"XA": "la Panadella",
	"XB": "la Llacuna",
	"XC": "Castellbisbal",
	"XD": "Ulldemolins",
	"XE": "Tarragona - Complex Educatiu",
	"XF": "Sabadell - Parc Agrari",
	"XG": "Parets del Vallès",
	"XH": "Sort",
	"XI": "Mollerussa",
	"XJ": "Girona",
	"XK": "Puig Sesolles (1.668 m)",
	"XL": "el Prat de Llobregat",
	"XM": "els Alamús",
	"XN": "Seròs",
	"XO": "Vic",
	"XP": "Gandesa",
	"XQ": "Tremp",
	"XR": "Prades",
	"XS": "Santa Coloma de Farners",
	"XT": "Solsona",
	"XU": "Canyelles",
	"XV": "Sant Cugat del Vallès - CAR",
	"XX": "Tornabous",
	"XY": "Alcarràs",
	"XZ": "Torroella de Fluvià",
	"X1": "Falset",
	"X2": "Barcelona - Zoo",
	"X3": "Alguaire",
	"X4": "Barcelona - el Raval",
	"X5": "PN dels Ports",
	"X6": "Baldomar",
	"X7": "Torres de Segre",
	"X8": "Barcelona - Zona Universitària",
	"X9": "Caldes de Montbui",
	"YA": "Puigcerdà",
	"YB": "Olot",
	"YC": "la Pobla de Segur",
	"YD": "les Borges Blanques",
	"YE": "Massoteres",
	"YF": "Mont-roig del Camp",
	"YG": "Tírvia",
	"YH": "Pujalt",
	"YJ": "Lleida - la Femosa",
	"YK": "Terrassa",
	"YL": "Riudecanyes",
	"YM": "Granollers",
	"Y4": "Alinyà",
	"Y5": "Navata",
	"Y6": "Tivissa",
	"Y7": "Port de Barcelona - Bocana Sud",
	"ZB": "Salòria (2.451 m)",
	"ZC": "Ulldeter (2.410 m)",
	"ZD": "la Tosa d'Alp 2500",
	"Z1": "Bonaigua (2.266 m)",
	"Z2": "Boí (2.535 m)",
	"Z3": "Malniu (2.230 m)",
	"Z5": "Certascan (2.400 m)",
	"Z6": "Sasseuva (2.228 m)",
	"Z7": "Espot (2.519 m)",
	"Z8": "el Port del Comte (2.316 m)",
	"Z9": "Cadí Nord (2.143 m) - Prat d'Aguiló",
}

// CodisVariables holds all the measurements performed by the weather stations. Note that not all the stations
// measures all the variables. To see which measurements a station does check the method MeasurementMetadataAllByStation
var CodisVariables = map[string]string{
	"1":  "Pressió atmosfèrica màxima",
	"2":  "Pressió atmosfèrica mínima",
	"3":  "Humitat relativa màxima",
	"X4": "Temperatura màxima de subsòl a 5 cm",
	"5":  "Temperatura mínima de subsòl a 5 cm",
	"6":  "TDR màxima a 10 cm",
	"7":  "TDR mínima a 10 cm",
	"8":  "Desviació estàndard de la irradiància neta",
	"9":  "Irradiància reflectida",
	"10": "Irradiància fotosintèticament activa (PAR)",
	"11": "Temperatura de supefície",
	"12": "Temperatura màxima de superfície",
	"13": "Temperatura mínima de superfície",
	"14": "Temperatura de subsòl a 40 cm",
	"16": "Nivell evaporímetre",
	"20": "Velocitat del vent a 10 m (vec.)",
	"21": "Direcció del vent a 10 m (m. u)",
	"22": "Desviació est. de la direcció del vent a 10 m",
	"23": "Velocitat del vent a 6 m (vec.)",
	"24": "Direcció del vent a 6 m (m. u)",
	"25": "Desviació est. de la direcció de vent a 6 m",
	"26": "Velocitat del vent a 2 m (vec.)",
	"27": "Direcció del vent a 2 m (m. u)",
	"28": "Desviació est. de la direcció del vent a 2 m",
	"30": "Velocitat del vent a 10 m (esc.)",
	"31": "Direcció de vent 10 m (m. 1)",
	"32": "Temperatura",
	"33": "Humitat relativa",
	"34": "Pressió atmosfèrica",
	"35": "Precipitació",
	"36": "Irradiància solar global",
	"37": "Desviació est. de la irradiància solar global",
	"38": "Gruix de neu a terra",
	"39": "Radiació UV",
	"40": "Temperatura màxima",
	"42": "Temperatura mínima",
	"44": "Humitat relativa mínima",
	"46": "Velocitat del vent a 2 m (esc.)",
	"47": "Direcció del vent a 2 m (m. 1)",
	"48": "Velocitat del vent a 6 m (esc.)",
	"49": "Direcció del vent a 6 m (m. 1)",
	"50": "Ratxa màxima del vent a 10 m",
	"51": "Direcció de la ratxa màxima del vent a 10 m",
	"53": "Ratxa màxima del vent a 6 m",
	"54": "Direcció de la ratxa màxima del vent a 6 m",
	"56": "Ratxa màxima del vent a 2 m",
	"57": "Direcció de la ratxa màxima del vent a 2 m",
	"59": "Irradiància neta",
	"60": "Temperatura de subsòl a 5 cm",
	"61": "Temperatura de subsòl a 50 cm",
	"62": "TDR a 10 cm",
	"63": "TDR a 35 cm",
	"64": "Humectació moll",
	"65": "Humectació sec",
	"66": "Humectació res",
	"67": "Humectació moll 2",
	"68": "Humectació sec 2",
	"69": "Humectació res 2",
	"70": "Precipitació acumulada",
	"71": "Bateria",
	"72": "Precipitació màxima en 1 minut",
	"74": "Humitat del combustible forestal 1",
	"75": "Temperatura del combustible forestal 1",
	"76": "Humitat del combustible forestal 2",
	"77": "Temperatura del combustible forestal 2",
	"78": "Humitat del combustible forestal 3",
	"79": "Temperatura del combustible forestal 3",
	"80": "Temperatura de la neu 1",
	"81": "Temperatura de la neu 2",
	"82": "Temperatura de la neu 3",
	"83": "Temperatura de la neu 4",
	"84": "Temperatura de la neu 5",
	"85": "Temperatura de la neu 6",
	"86": "Temperatura de la neu 7",
	"87": "Temperatura de la neu 8",
	"88": "Quality number",
	"89": "Temperatura del datalogger",
	"90": "Altura màxima",
	"91": "Període màxima",
	"92": "Altura significant",
	"93": "Període significant",
	"94": "Altura mitjana",
	"95": "Període mitjà",
	"96": "Direcció del pic",
	"97": "Temperatura superficial del mar",
}

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

// Parameters holds all the options to be passed in to the methods
type Parameters struct {
	codiEstacio  string // should reference a key in the CodisEstacions map
	codiVariable string // should reference a key in the CodisVariables map
	codiEstat    string // should reference a key in the CodisEstat map
	Data
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

// APIError returned on failed API calls.
type APIError struct {
	Message string `json:"message"`
	COD     string `json:"cod"`
}

// Variables to calculate time values in the trace request executions
var t0, t1, t2, t3, t4, t5, t6 time.Time

// Wraper of Fprintf to colorize the output
func printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(color.Output, format, a...)
}

// Color settings
func grayscale(code color.Attribute) func(string, ...interface{}) string {
	return color.New(code + 232).SprintfFunc()
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
	trace  *httptrace.ClientTrace

	//cr *resty.Client
}

// NewSettings returns a new Setting pointer with default http client.
func NewSettings() *Settings {
	return &Settings{
		client: http.DefaultClient,
	}
}

// Optional client settings
type Option func(s *Settings) error

//// WithHttpClient sets custom http client when creating a new Client.
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
