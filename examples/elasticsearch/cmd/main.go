package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/oscaromeu/meteocat"
	"github.com/oscaromeu/meteocat/examples/elasticsearch/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Entry point of the app

var (
	log         = logrus.New()
	appdir      = os.Getenv("PWD")
	homeDir     string
	pidFile     string
	logDir      = filepath.Join(appdir, "log")
	confDir     = filepath.Join(appdir, "config")
	downloadDir = filepath.Join(appdir, "download")
	dataDir     = confDir
	configFile  = filepath.Join(confDir, "conf.yaml")
	// MainConfig has all configuration
	MainConfig config.Config
)

func writePIDFile() {
	if pidFile == "" {
		return
	}

	// Ensure the required directory structure exists.
	err := os.MkdirAll(filepath.Dir(pidFile), 0700)
	if err != nil {
		log.Fatal(3, "Failed to verify pid directory", err)
	}

	// Retrieve the PID and write it.
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(pidFile, []byte(pid), 0644); err != nil {
		log.Fatal(3, "Failed to write pidfile", err)
	}
}

func init() {

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.Formatter = customFormatter
	customFormatter.FullTimestamp = true

	//log.SetOutput(os.Stdout)
	v := viper.New()

	// now load up config settings
	if _, err := os.Stat(configFile); err == nil {
		v.SetConfigFile(configFile)
		confDir = filepath.Dir(configFile)
	} else {
		v.SetConfigName("conf")
		v.SetConfigType("yaml")
		v.AddConfigPath("../config")
		v.AddConfigPath("../../config")
		v.AddConfigPath(".")
	}
	err := v.ReadInConfig()
	if err != nil {
		log.Errorf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}
	err = v.Unmarshal(&MainConfig)
	if err != nil {
		log.Errorf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}
	log.Infof("CONFIG FROM FILE : %+v", &MainConfig)

	//parse config with environtment vars

	//envConf := &config.Config{}

	err = envconfig.Process("M_", &MainConfig)
	if err != nil {
		log.Warnf("Some error happened when trying to read config from env: %s", err)
	}
	//log.Infof("CONFIG FROM ENV : %+v", envConf)

	//mergo.MergeWithOverwrite(&agent.MainConfig, envConf)

	log.Infof("CONFIG AFTER MERGE : %+v", &MainConfig)
	// Setting defaults

	cfg := &MainConfig

	if cfg.General.LogMode != "file" {
		//default if not set
		log.Out = os.Stdout

	} else {
		if len(cfg.General.LogDir) > 0 {
			logDir = cfg.General.LogDir
			os.Mkdir(logDir, 0755)
			//Log output
			f, _ := os.OpenFile(logDir+"/es_meteocat.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
			log.Out = f
		}
	}
	if len(cfg.General.LogLevel) > 0 {
		l, _ := logrus.ParseLevel(cfg.General.LogLevel)
		log.Level = l
	}

	if len(cfg.General.HomeDir) > 0 {
		homeDir = cfg.General.HomeDir
	}

	log.Debugf("MAINCONFIG LOAD  %#+v", cfg)
	fmt.Println("")

}

func main() {

	// Load Logger settings for other packages. These packages will inherit the log config defined in main
	config.SetLogger(log)

	d, err := meteocat.NewMesurades(MainConfig.General.ApiKey)
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
		meteocat.OptionCodiVariable("32"),
		meteocat.OptionData(data),
	)

	// Call MeasurementByDay Method
	d.MeasurementLast(params)

	for _,v := range d.Measurements {
		fmt.Println(v)
	}
}
