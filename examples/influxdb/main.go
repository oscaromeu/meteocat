package main

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	client "github.com/influxdata/influxdb-client-go/v2"
	"github.com/oscaromeu/meteocat"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

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

var (
	log = logrus.New()
	cfg Config
	t time.Time

	Fields = make(map[string]interface{})
	Tags = make(map[string]string)
)

const layout = "2006-01-02T15:04Z"

type General struct {
	ApiKey string `yaml:"meteocat_api_key" env:"METEOCAT_API_KEY" env-description:"Meteocat Api Key"`
}

type Database struct {
	InfluxServer      string `yaml:"influx_server" env:"INFLUX_SERVER" env-description:"Influx (Host) Server Instance"`
	InfluxPort        string `yaml:"influx_port" env:"INFLUX_PORT" env-description:"Influx Server Instance"`
	InfluxBucket      string `yaml:"influx_bucket" env:"INFLUX_BUCKET" env-description:"Influx DB Instance"`
	InfluxMeasurement string `yaml:"influx_measurement" env:"INFLUX_MEASUREMENT" default:"" env-description:"Influx Measurement Name"` //TODO: Name influx measurement
	InfluxOrg         string `yaml:"influx_org" env:"INFLUX_ORG" env-description:"Influx Username of Server Instance"`
	InfluxToken       string `yaml:"influx_token" env:"INFLUX_TOKEN" env-description:"Influx Password of Server Instance"`
	//	ExtraTags            map[string]string `yaml:"extra_tags" env:"EXTRA_TAGS" env-description:"Extra tags name to add to the measurements"`
	LogDebug bool `yaml:"log_debug" env:"LOG_DEBUG" default:"false" env-description:"Log Level"`
}

// Config is a application configuration structure
type Config struct {
	General
	Database
}

// Args command-line parameters
type Args struct {
	ConfigPath string
}


func init() {

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.Formatter = customFormatter
	customFormatter.FullTimestamp = true

	log.SetOutput(os.Stdout)

	args := ProcessArgs(&cfg)

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, &cfg); err != nil {
		log.Fatal(err.Error())
		os.Exit(2)
	}

	if cfg.LogDebug {

		log.SetLevel(logrus.DebugLevel)
	}

	log.Debugf("%#+v", cfg.Database)

}

// ProcessArgs processes and handles CLI arguments
func ProcessArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("Configuration file", 1)
	f.StringVar(&a.ConfigPath, "c", "config.yml", "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}

func InitDB(conf Config) (client.Client, error) {

	fmt.Sprintf("http://%s:%s", conf.InfluxServer, conf.InfluxPort)

	c := client.NewClient(fmt.Sprintf("http://%s:%s", conf.InfluxServer, conf.InfluxPort), conf.InfluxToken)

	return c, nil

}

func main() {


	c, err := InitDB(cfg)
	if err != nil {
		log.Infof("ERROR: %s", err)
		os.Exit(1)
	}

	d, err := meteocat.NewMesurades(cfg.ApiKey)
	if err != nil {
		log.Fatalln(err)
		os.Exit(3)
	}

	if meteocat.CheckAPIKeyExists(d.Key) == false {
		fmt.Println("ApiKey is not set. ")
	}


	year, month, day := time.Now().Date()
	data := meteocat.Data{
		Any: strconv.Itoa(year),
		Mes: "0" + strconv.Itoa(int(month)),
		Dia: strconv.Itoa(day - 1),
	}

	params, _ := meteocat.NewParameters(
		meteocat.OptionCodiEstacio("UG"),
		//meteocat.OptionCodiVariable("5"),
		meteocat.OptionData(data),
	)

	d.MeasurementAllByStation(params)

	writeAPI := c.WriteAPI(cfg.InfluxOrg, cfg.InfluxBucket)
	for _, v := range d.Measurements {
		Tags["codi_estacio"] = v.Codi
		Tags["Estacio"] = CodisEstacions[v.Codi]
		for _, lec := range v.Variables {
			Tags["codi_variable"] = strconv.Itoa(lec.Codi)
			Tags["Variable"] = CodisVariables[strconv.Itoa(lec.Codi)]
			for _, l := range lec.Lectures {
				t, _ = time.Parse(layout, l.Data)
				Fields["valor"] = l.Valor

				p := client.NewPoint(cfg.InfluxMeasurement,
					Tags,
					Fields,
					t)
				// write point asynchronously
				writeAPI.WritePoint(p)
			}
		}
		// Flush writes
		writeAPI.Flush()
	}

	// always close client at the end
	defer c.Close()

}