package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var (

	//Log the Logger
	log *logrus.Logger
)



// ReadEnv reads the environmental variable associated with a given key function paramater
func ReadEnv(key string) string {

	v := viper.New()

	// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use.
	// E.g. In this case the prefix is "gs", the env registry will look for env
	// variables that start with "GS_". To import a template to a zabbix instance we must set the
	// following variables: ZABBIX_PASSWORD, ZABBIX_SERVER, ZABBIX_USERNAME
	v.SetEnvPrefix("gs")
	//v.SetDefault("LOG_DEBUG", false)

	// Make sure that ENV.VAR => ENV_VAR
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Enable VIPER to read Environment Variables
	v.AutomaticEnv()

	// AllowEmptyEnv tells Viper to consider set,
	// but empty environment variables as valid values instead of falling back.
	v.AllowEmptyEnv(true)

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value, err := v.Get(key).(string)

	if !err {
		err := errors.Errorf("Failed to get the environmental variable %q", value)
		log.Error(err)
	}
	return value
}

//// Load reads the config/config.yaml file
//func Load() (*Config, error) {
//
//	v := viper.New()
//	// Set the file name of the configurations file
//	v.SetConfigName("conf")
//
//	// Set the path to look for the configurations file
//	v.AddConfigPath("../../config")
//
//	// Enable VIPER to read Environment Variables
//	v.AutomaticEnv()
//
//	v.SetConfigType("yaml")
//	//var configuration Configurations
//
//	var conf Config
//
//	if err := v.ReadInConfig(); err != nil {
//		return nil, errors.Wrapf(err, "Error reading config file")
//	}
//
//	err := v.Unmarshal(&conf)
//	if err != nil {
//		return nil, errors.Wrapf(err, "Unable to decode into struct, %q", conf)
//	}
//	return &conf, err
//
//}

// SetLogger set the output log
func SetLogger(l *logrus.Logger) {
	log = l
}
