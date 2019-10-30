package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"gitlab.com/nl2go/hrobot-cli/cmd"
	"gitlab.com/nl2go/hrobot-cli/config"
)

func main() {
	log.SetOutput(os.Stdout)

	var cfg config.Config
	err := envconfig.Process("hrobotcli", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	rootCmd := cmd.NewRootCommand(log.StandardLogger(), &cfg)
	rootCmd.SetErr(log.StandardLogger().Out)
	if err := rootCmd.Execute(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}
