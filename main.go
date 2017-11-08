package main

import (
	"flag"
	"fmt"
	"log"
	"playing-with-refelction-in-go/config"
)

func main() {
	//Make sure the file exists
	configPath := flag.String("config", "config/config.json", "path of the config file")

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.BaseURL)
	fmt.Println(cfg.Name)
	fmt.Println(cfg.Driver)

	err = cfg.UseCustomEnvConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.BaseURL)
	fmt.Println(cfg.Name)
	fmt.Println(cfg.Driver)
}
