package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	client "github.com/influxdata/influxdb-client-go/v2"
	"github.com/oscaromeu/meteocat"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
	cfg Config
	t   time.Time

	Fields = make(map[string]interface{})
	Tags   = make(map[string]string)
)

type General struct {
	ApiKey  string   `yaml:"meteocat_api_key" env:"METEOCAT_API_KEY" env-description:"Meteocat Api Key"`
	Estacio []string `yaml:"meteocat_codi_estacio" env:"METEOCAT_CODI_ESTACIO" env-description:"Meteocat Codi Estacio"`
}

const layout = `2006-01-02T15:04:05.000`

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

func timeShift(now time.Time, timeType string, shift int) time.Time {
	var shiftTime time.Time
	switch timeType {
	case "year":
		shiftTime = now.AddDate(shift, 0, 0)
	case "month":
		shiftTime = now.AddDate(0, shift, 0)
	case "day":
		shiftTime = now.AddDate(0, 0, shift)
	case "hour":
		shiftTime = now.Add(time.Hour * time.Duration(shift))
	case "minute":
		shiftTime = now.Add(time.Minute * time.Duration(shift))
	case "second":
		shiftTime = now.Add(time.Second * time.Duration(shift))
	default:
		shiftTime = now
	}
	return shiftTime
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

	c := client.NewClient("http://influxdb:8086", fmt.Sprintf("%s:%s", "", ""))
	return c, nil
}

func main() {

	c, err := InitDB(cfg)
	if err != nil {
		log.Infof("ERROR: %s", err)
		os.Exit(1)
	}

	// Shift time since opendata has a delay of two hours
	tt := timeShift(time.Now(), "hour", -2)

	d, _ := meteocat.NewOpenDataMesurades()

	data := meteocat.Data{
		Any: strconv.Itoa(tt.Year()),
		Mes: strconv.Itoa(int(tt.Month())),
		Dia: strconv.Itoa(tt.Day()),
	}

	// Get the current minute and round to 30 or 00
	// Opendata precision is rouglhy 30min
	var minute string
	if (tt.Minute() >= 30) && (tt.Minute() < 45) {
		minute = "30"
	} else {
		minute = "00"
	}

	timeDate := meteocat.TimeDate{
		Hour:    strconv.Itoa(tt.Hour()),
		Minute:  minute,
		Seconds: "00",
	}

	fmt.Println(timeDate)
	p, _ := meteocat.NewParameters(
		meteocat.OptionData(data),
		meteocat.OptionTimeDate(timeDate),
	)

	d.OpenDataMeasurementAllByStation(p)
	log.Info(d.OpenData)

	writeAPI := c.WriteAPI("", "meteocat/autogen")

	for _, v := range d.OpenData {
		meas, err := strconv.ParseFloat(v.ValorLectura, 64)
		if err != nil {
			log.Infof("ERROR: %s", err)
		}
		// Extract and parse time from the results
		// The point stored in influx will have the time from opendata results
		zz, err := time.Parse(layout, v.DataLectura)
		// create point using fluent style
		p := client.NewPointWithMeasurement(cfg.InfluxMeasurement).
			AddTag("codi_estacio", v.CodiEstacio).
			AddTag("nom_estacio", meteocat.CodisEstacions[v.CodiEstacio]).
			AddTag("codi_variable", v.CodiVariable).
			AddTag("nom_variable", meteocat.CodisVariables[v.CodiVariable]).
			AddField("measurement", meas).
			SetTime(zz)
		writeAPI.WritePoint(p)
		// Flush writes
		writeAPI.Flush()

	}
	// always close client at the end
	defer c.Close()

}
