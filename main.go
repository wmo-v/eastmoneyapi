package main

import (
	"eastmoneyapi/config"
	"eastmoneyapi/service"
	"flag"
	"log"
	math_rand "math/rand"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var configPath string

func init() {
	math_rand.Seed(time.Now().Unix())
	log.SetFlags(log.Lshortfile)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
		ForceColors:     true,
	})
}
func init() {
	flag.StringVar(&configPath, "config", "", "")
	flag.Parse()
	if configPath != "" {
		config.SetConfigPath(configPath)
	}
}

// 每周1-5 9.32
var openTimeSpec = "32 9 * * 1-5"

// 每周1-5 15.00
var closeTimeSpec = "0 15 * * 1-5"

func main() {
	c := cron.New()
	c.AddFunc(openTimeSpec, func() {
		z := service.NewZ513050Svc()
		z.Start()
		var id cron.EntryID
		id, _ = c.AddFunc(closeTimeSpec, func() {
			z.Close()
			c.Remove(id)
		})
	})
	c.Run()
}
