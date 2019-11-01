package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

func (app *RobotApp) NewIPGetListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ip:list",
		Short: "Print list of IP's",
		Long:  "Print list of IP's in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			ips, err := app.client.IPGetList()
			if err != nil {
				app.logger.Errorln(err)
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
