package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	client "github.com/influxdata/influxdb-client-go/v2"
	"github.com/oscaromeu/meteocat"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	log = logrus.New()
	cfg Config
	t   time.Time

	Fields = make(map[string]interface{})
	Tags   = make(map[string]string)
)

var (
	timeApiBaseURL = "https://www.timeapi.io/api/Time/current/zone?timeZone=Europe/Madrid"
)

type Time struct {
	Year         int    `json:"year"`
	Month        int    `json:"month"`
	Day          int    `json:"day"`
	Hour         int    `json:"hour"`
	Minute       int    `json:"minute"`
	Seconds      int    `json:"seconds"`
	MilliSeconds int    `json:"milliSeconds"`
	DateTime     string `json:"dateTime"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	TimeZone     string `json:"timeZone"`
	DayOfWeek    string `json:"dayOfWeek"`
	DstActive    bool   `json:"dstActive"`
}

type General struct {
	ApiKey  string   `yaml:"meteocat_api_key" env:"METEOCAT_API_KEY" env-description:"Meteocat Api Key"`
	Estacio []string `yaml:"meteocat_codi_estacio" env:"METEOCAT_CODI_ESTACIO" env-description:"Meteocat Codi Estacio"`
}

const layout = "2006-01-02T15:04Z"

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

func (t *Time) getTime() {

	url := timeApiBaseURL
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Println(err)
		return
	}

}

func main() {

	c, err := InitDB(cfg)
	if err != nil {
		log.Infof("ERROR: %s", err)
		os.Exit(1)
	}

	var t Time
	var Tags = make(map[string]string)

	t.getTime()

	d, _ := meteocat.NewOpenDataMesurades()

	data := meteocat.Data{
		Any: strconv.Itoa(t.Year),
		Mes: strconv.Itoa(t.Month),
		Dia: strconv.Itoa(t.Day),
	}

	//	timeDate := meteocat.TimeDate{
	//		Hour:    strconv.Itoa(t.Hour),
	//		Minute:  strconv.Itoa(t.Minute),
	//		Seconds: "00",
	//	}

	timeDate := meteocat.TimeDate{
		Hour:    "05",
		Minute:  "00",
		Seconds: "00",
	}
	fmt.Println(data)
	p, _ := meteocat.NewParameters(
		meteocat.OptionData(data),
		meteocat.OptionTimeDate(timeDate),
	)

	d.OpenDataMeasurementAllByStation(p)

	writeAPI := c.WriteAPI(cfg.InfluxOrg, cfg.InfluxBucket)
	for _, v := range d.OpenData {
		Tags["codi_estacio"] = v.CodiEstacio
		Tags["nom_estacio"] = meteocat.CodisEstacions[v.CodiEstacio]
		Tags["codi_variable"] = v.CodiVariable
		Tags["nom_variable"] = meteocat.CodisVariables[v.CodiVariable]
		t, _ = time.Parse(layout, data.Any, data.Mes, data.Dia)
		Fields["valor"] = l.Valor
		p := client.NewPoint(cfg.InfluxMeasurement,
			Tags,
			Fields, fmt.Sprintf("?data_lectura=%s-%s-%sT%s:%s:%s", p.Any, p.Mes, p.Dia, p.Hour, p.Minute, p.Seconds))
		// write point asynchronously
		writeAPI.WritePoint(p)
	}
	// Flush writes
	writeAPI.Flush()
	// always close client at the end
	defer c.Close()

}
