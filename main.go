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
	var c = client.NewEastMoneyClient()
	if err := c.Login(config.GetConfg().User.Account, config.GetConfg().User.Password); err != nil {
		panic(err)
	}
	// 根据自己的需求进行交易
}
