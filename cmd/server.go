package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.com/nl2go/hrobot-cli/client"
	"gitlab.com/nl2go/hrobot-cli/config"
)

func NewServerCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server:get-list",
		Short: "Print list of servers",
		Long:  "Print list of servers in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
			servers, err := robotClient.ServerGetList()

			if err != nil {
				logger.Errorln(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"id", "ip", "name", "datacenter"})

			for _, server := range *servers {
				t.AppendRow(table.Row{
					server.Server.ServerNumber,
					server.Server.ServerIP,
					server.Server.ServerName,
					server.Server.Dc,
				})
			}

			t.AppendFooter(table.Row{"", "", "Total", len(*servers)})
			t.Render()
		},
	}
}
