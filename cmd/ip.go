package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.com/nl2go/hrobot-cli/client"
	"gitlab.com/nl2go/hrobot-cli/config"
)

func NewIPGetListCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "ip:list",
		Short: "Print list of IP's",
		Long:  "Print list of IP's in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
			ips, err := robotClient.IPGetList()
			if err != nil {
				logger.Errorln(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"ip", "server_ip", "server_number", "locked"})

			for _, ip := range ips {
				t.AppendRow(table.Row{
					ip.IP,
					ip.ServerIP,
					ip.ServerNumber,
					ip.Locked,
				})
			}

			t.SortBy([]table.SortBy{
				table.SortBy{
					Name: "server_ip",
					Mode: table.Asc,
				},
			})
			t.AppendFooter(table.Row{"", "", "Total", len(ips)})
			t.Render()
		},
	}
}
