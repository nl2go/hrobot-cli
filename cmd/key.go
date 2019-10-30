package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.com/nl2go/hrobot-cli/client"
	"gitlab.com/nl2go/hrobot-cli/config"
)

func NewKeyCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "key:get-list",
		Short: "Print list of ssh keys",
		Long:  "Print list of ssh keys in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
			keys, err := robotClient.KeyGetList()

			if err != nil {
				logger.Errorln(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"name", "type", "size", "fingerprint"})

			for _, key := range *keys {
				t.AppendRow(table.Row{
					key.Key.Name,
					key.Key.Type,
					key.Key.Size,
					key.Key.Fingerprint,
				})
			}

			t.AppendFooter(table.Row{"", "", "Total", len(*keys)})
			t.Render()
		},
	}
}
