package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"gitlab.com/nl2go/hrobot-cli/config"
)

func NewRootCommand(logger *log.Logger, cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hrobot-cli",
		Short: "CLI application for the hetzner robot API",
	}

	rootCmd.AddCommand(NewServerCmd(logger, cfg))

	return rootCmd
}
