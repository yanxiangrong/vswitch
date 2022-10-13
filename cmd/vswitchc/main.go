package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"vswitch/client"
	"vswitch/pkg/config"
	"vswitch/pkg/util/log"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cfg := config.ClientConf{
		ServerAddr:      "hongkong.yandage.top",
		ServerPort:      67,
		LogFile:         "console",
		LogWay:          "console",
		LogLevel:        "debug",
		LogMaxDays:      3,
		DisableLogColor: false,
	}
	log.InitLog(cfg.LogWay, cfg.LogFile, cfg.LogLevel, cfg.LogMaxDays, cfg.DisableLogColor)
	log.Debug(fmt.Sprint("Log laval: ", cfg.LogLevel))

	svr, err := client.NewService(cfg)
	if err != nil {
		log.Error(fmt.Sprint(err))
		os.Exit(-1)
	}
	log.Info("vswitchc started successfully")
	svr.Run()
}
