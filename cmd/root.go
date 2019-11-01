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

	rootCmd.AddCommand(NewServerGetListCmd(logger, cfg))
	rootCmd.AddCommand(NewServerGetCmd(logger, cfg))
	rootCmd.AddCommand(NewServerReversalCmd(logger, cfg))
	rootCmd.AddCommand(NewServerSetNameCmd(logger, cfg))
	rootCmd.AddCommand(NewServerGenerateAnsibleInventoryCmd(logger, cfg))
	rootCmd.AddCommand(NewKeyGetListCmd(logger, cfg))
	rootCmd.AddCommand(NewIPGetListCmd(logger, cfg))
	rootCmd.AddCommand(NewRdnsGetListCmd(logger, cfg))
	rootCmd.AddCommand(NewRdnsGetCmd(logger, cfg))

	return rootCmd
}
