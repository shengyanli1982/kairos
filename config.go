package kairos

type Config struct {
	callback Callback
}

func NewConfig() *Config {
	return &Config{}
}

func isConfigValid(conf *Config) *Config {
	return conf
}
