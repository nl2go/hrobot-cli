package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/nl2go/hrobot-cli/cmd"
	"github.com/nl2go/hrobot-cli/config"
	client "github.com/nl2go/hrobot-go"
)

func main() {
	log.SetOutput(os.Stdout)

	var cfg config.Config
	err := envconfig.Process("hrobotcli", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
	hrobotApp := cmd.NewRobotApp(robotClient, log.StandardLogger())
	if err := hrobotApp.Run(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}
