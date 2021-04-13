package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	client "github.com/nl2go/hrobot-go"
)

const version = "0.1.1"
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
		Short: fmt.Sprintf("CLI application for the Hetzner Robot API - version %s", version),
	}

	rootCmd.AddCommand(app.NewServerGetListCmd())
	rootCmd.AddCommand(app.NewServerGetCmd())
	rootCmd.AddCommand(app.NewServerReversalCmd())
	rootCmd.AddCommand(app.NewServerSetNameCmd())
	rootCmd.AddCommand(app.NewServerActivateRescueCmd())
	rootCmd.AddCommand(app.NewServerResetCmd())
	rootCmd.AddCommand(app.NewServerGenerateAnsibleInventoryCmd())
	rootCmd.AddCommand(app.NewKeyGetListCmd())
	rootCmd.AddCommand(app.NewIPGetListCmd())
	rootCmd.AddCommand(app.NewRdnsGetListCmd())
	rootCmd.AddCommand(app.NewRdnsGetCmd())
	rootCmd.AddCommand(app.NewFailoverGetListCmd())
	rootCmd.AddCommand(app.NewFailoverGetCmd())
	rootCmd.AddCommand(app.NewVersionCmd())

	return rootCmd
}

func (app *RobotApp) NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of hrobot-cli",
		Long:  `All software has versions. This is hrobot-cli's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Hetzner Robot Webservice command line interface version: %s", version))
		},
	}
}
