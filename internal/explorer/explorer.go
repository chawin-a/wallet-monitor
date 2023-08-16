package explorer

type Config struct {
	ApiURL string `mapstructure:"api_url"`
	ApiKey string `mapstructure:"api_key"`
}

type Explorer struct {
	ApiURL string
	ApiKey string
}

func NewExplorer(conf *Config) *Explorer {
	return &Explorer{
		ApiURL: conf.ApiURL,
		ApiKey: conf.ApiKey,
	}
}
