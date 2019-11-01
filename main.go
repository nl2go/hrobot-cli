package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"gitlab.com/newsletter2go/hrobot-cli/cmd"
	"gitlab.com/newsletter2go/hrobot-cli/config"
)

func main() {
	log.SetOutput(os.Stdout)

	var cfg config.Config
	err := envconfig.Process("hrobotcli", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	hrobotApp := cmd.NewRobotApp(log.StandardLogger(), &cfg)
	if err := hrobotApp.Run(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}
