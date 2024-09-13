package main

import (
	"flag"
	"log"

	"github.com/konocha/qr-generate/internal/app/apiserver"
)

var(
	configPath string
)

func init(){
	flag.StringVar(&configPath, "config-path", "configs/qrgenerate.toml", "path to config file")
}

func main() {
	flag.Parse()

	cfg := apiserver.NewConfig()

	err := apiserver.Start(cfg)
	if err != nil{
		log.Fatal(err)
	}
}