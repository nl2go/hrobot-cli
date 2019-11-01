package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"gitlab.com/newsletter2go/hrobot-go"
	"gitlab.com/newsletter2go/hrobot-cli/config"
)

const version = "1.0.0"
const userAgent = "hrobot-cli/" + version

type RobotApp struct {
	logger *log.Logger
	cfg    *config.Config
	client client.RobotClient
}

func NewRobotApp(logger *log.Logger, cfg *config.Config) *RobotApp {
	robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
	robotClient.SetUserAgent(userAgent)

	return &RobotApp{
		logger: logger,
		cfg:    cfg,
		client: robotClient,
	}
}

func (app *RobotApp) Run() error {
	rootCmd := app.NewRootCommand(app.logger, app.cfg)
	rootCmd.SetErr(app.logger.Out)

	err := rootCmd.Execute()
	return err
}

func (app *RobotApp) NewRootCommand(logger *log.Logger, cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hrobot-cli",
		Short: "CLI application for the hetzner robot API",
	}

	rootCmd.AddCommand(app.NewServerGetListCmd())
	rootCmd.AddCommand(app.NewServerGetCmd())
	rootCmd.AddCommand(app.NewServerReversalCmd())
	rootCmd.AddCommand(app.NewServerSetNameCmd())
	rootCmd.AddCommand(app.NewServerGenerateAnsibleInventoryCmd())
	rootCmd.AddCommand(app.NewKeyGetListCmd())
	rootCmd.AddCommand(app.NewIPGetListCmd())
	rootCmd.AddCommand(app.NewRdnsGetListCmd())
	rootCmd.AddCommand(app.NewRdnsGetCmd())

	return rootCmd
}
