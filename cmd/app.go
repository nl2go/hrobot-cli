package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	client "gitlab.com/newsletter2go/hrobot-go"
)

const version = "0.0.0"
const userAgent = "hrobot-cli/" + version

type RobotApp struct {
	logger *log.Logger
	client client.RobotClient
}

func NewRobotApp(robotClient client.RobotClient, logger *log.Logger) *RobotApp {
	robotClient.SetUserAgent(userAgent)

	return &RobotApp{
		logger: logger,
		client: robotClient,
	}
}

func (app *RobotApp) Run() error {
	rootCmd := app.NewRootCommand(app.logger)
	rootCmd.SetErr(app.logger.Out)

	err := rootCmd.Execute()
	return err
}

func (app *RobotApp) NewRootCommand(logger *log.Logger) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hrobot-cli",
		Short: "CLI application for the hetzner robot API",
	}

	rootCmd.AddCommand(app.NewServerGetListCmd())
	rootCmd.AddCommand(app.NewServerGetCmd())
	rootCmd.AddCommand(app.NewServerReversalCmd())
	rootCmd.AddCommand(app.NewServerSetNameCmd())
	rootCmd.AddCommand(app.NewServerActivateRescueCmd())
	rootCmd.AddCommand(app.NewServerGenerateAnsibleInventoryCmd())
	rootCmd.AddCommand(app.NewKeyGetListCmd())
	rootCmd.AddCommand(app.NewIPGetListCmd())
	rootCmd.AddCommand(app.NewRdnsGetListCmd())
	rootCmd.AddCommand(app.NewRdnsGetCmd())

	return rootCmd
}
