package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/nl2go/hrobot-go/models"
)

func (app *RobotApp) NewFailoverGetListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "failover:list",
		Short: "Print list of failover IP's",
		Long:  "Print list of failover IP's in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			failoverIPList, err := app.client.FailoverGetList()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"ip", "server number", "active server IP"})

			for _, failoverIP := range failoverIPList {
				t.AppendRow(table.Row{
					failoverIP.IP,
					failoverIP.ServerNumber,
					failoverIP.ActiveServerIP,
				})
			}

			t.AppendFooter(table.Row{"Total", len(failoverIPList)})
			t.Render()
		},
	}
}

func (app *RobotApp) NewFailoverGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "failover:get",
		Short: "Print single failover IP",
		Long: `Print details of single failover IP in hetzner account
		failover IP can be chosen interactively`,
		Run: func(cmd *cobra.Command, args []string) {
			failoverIPList, err := app.client.FailoverGetList()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			prompt := promptui.Select{
				Label:             "Select failover IP",
				Items:             failoverIPList,
				Searcher:          getFailoverSearcher(failoverIPList),
				Size:              10,
				Templates:         getFailoverSelectTemplates(),
				StartInSearchMode: true,
			}

			choosenIdx, _, err := prompt.Run()
			if err != nil {
				app.logger.Errorln("Prompt failed: ", err)
				return
			}

			choosenFailover := failoverIPList[choosenIdx]
			fmt.Println("Chosen failover IP: ", choosenFailover.IP)

			// directly print info without additional get as that does not deliver more data
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"field", "value"})
			t.AppendRow(table.Row{"ip", choosenFailover.IP})
			t.AppendRow(table.Row{"net mask", choosenFailover.Netmask})
			t.AppendRow(table.Row{"server number", choosenFailover.ServerNumber})
			t.AppendRow(table.Row{"server ip", choosenFailover.ServerIP})
			t.AppendRow(table.Row{"active server IP", choosenFailover.ActiveServerIP})

			t.Render()
		},
	}
}

func getFailoverSearcher(failoverIPList []models.Failover) func(string, int) bool {
	return func(input string, index int) bool {
		failoverIP := failoverIPList[index]
		ip := strings.Replace(strings.ToLower(failoverIP.IP), " ", "", -1)
		activeServerIP := strings.Replace(strings.ToLower(failoverIP.ActiveServerIP), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(ip, input) || strings.Contains(activeServerIP, input)
	}
}

func getFailoverSelectTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "→ {{ .IP | green }} ({{ .ActiveServerIP | yellow }})",
		Inactive: "  {{ .IP | cyan }} ({{ .ActiveServerIP | red }})",
		Selected: "→ {{ .IP | cyan }}",
		Details: `
	--------- Selected server ----------
	{{ "IP:" | faint }}	          {{ .IP }}
	{{ "Active server IP:" | faint }}	  {{ .ActiveServerIP }}`,
	}
}
