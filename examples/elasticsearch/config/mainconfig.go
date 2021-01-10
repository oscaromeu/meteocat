package config

// GeneralConfig has miscellaneous configuration options
type GeneralConfig struct {
	LogDir         string `mapstructure:"log_dir" envconfig:"M_LOGDIR"`
	LogLevel       string `mapstructure:"log_level" envconfig:"M_LOGLEVEL"`
	LogMode        string `mapstructure:"log_mode" envconfig:"M_LOGMODE"`
	HomeDir        string `mapstructure:"home_dir" envconfig:"M_HOMEDIR"`
	TemporalPath   string `mapstructure:"temporal_path" envconfig:"M_TEMPORAL_PATH"`
	ESUser     string `mapstructure:"user" envconfig:"M_ES_USERNAME"`
	ESPassword string `mapstructure:"password" envconfig:"M_ES_PASSWORD"`
	ESServer      string `mapstructure:"password" envconfig:"M_ES_SERVER"`
	ApiKey      string `mapstructure:"api_key" envconfig:"M_API_KEY"`
}


type Config struct {
	General    GeneralConfig `mapstructure:"general"`
}
