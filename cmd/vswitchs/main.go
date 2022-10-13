package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"vswitch/pkg/config"
	"vswitch/pkg/util/log"
	"vswitch/server"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cfg := config.ServerConf{
		BindPort:        67,
		LogFile:         "console",
		LogWay:          "console",
		LogLevel:        "debug",
		LogMaxDays:      3,
		DisableLogColor: false,
	}
	log.InitLog(cfg.LogWay, cfg.LogFile, cfg.LogLevel, cfg.LogMaxDays, cfg.DisableLogColor)
	log.Debug(fmt.Sprint("Log laval: ", cfg.LogLevel))

	svr, err := server.NewService(cfg)
	if err != nil {
		log.Error(fmt.Sprint(err))
		os.Exit(-1)
	}
	log.Info("vswitchs started successfully")
	svr.Run()
}
