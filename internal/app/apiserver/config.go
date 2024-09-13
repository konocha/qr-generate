package apiserver

type Config struct{
	BindAddr string
	DataBaseURL string
}

func NewConfig() *Config{
	return &Config{
		BindAddr: ":8081",
		DataBaseURL: "root:Sofa=22082014@tcp(0.0.0.0:3306)/testDatabase",
	}
}