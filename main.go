package main

import (
	"eastmoneyapi/client"
	"eastmoneyapi/config"
	"flag"
	"log"
	math_rand "math/rand"
	"time"
)

var configPath string

func init() {
	math_rand.Seed(time.Now().Unix())
	log.SetFlags(log.Lshortfile)
}
func init() {
	flag.StringVar(&configPath, "config", "", "")
	flag.Parse()
	if configPath != "" {
		config.SetConfigPath(configPath)
	}
}

func main() {
	// z := service.NewZ513050Svc()
	// z.Start()
	c := client.NewEastMoneyClient()
	c.GetStockList()
}
