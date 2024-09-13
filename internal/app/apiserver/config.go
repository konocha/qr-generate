package apiserver

type Config struct{
	BindAddress string
	DataBaseURL string
}

func NewConfig(c Config) *Config{
	return &Config{
		BindAddress: c.BindAddress,
		DataBaseURL: c.DataBaseURL,
	}
}