package main

import (
	"flag"
	"log"

	"github.com/konocha/qr-generate/internal/app/apiserver"
	"github.com/BurntSushi/toml"
)

var(
	configPath string
)

func init(){
	flag.StringVar(&configPath, "config-path", "configs/qrgenerate.toml", "path to config file")
}



func main() {
	var config apiserver.Config
	flag.Parse()

	_, err := toml.DecodeFile(configPath, &config)
	if err != nil{
		log.Fatal()
	}
	cfg := apiserver.NewConfig(config)

	err = apiserver.Start(cfg)
	if err != nil{
		log.Fatal(err)
	}
}