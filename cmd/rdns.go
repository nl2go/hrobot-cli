package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"gitlab.com/newsletter2go/hrobot-go/models"
)

func (app *RobotApp) NewRdnsGetListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rdns:list",
		Short: "Print list of reverse DNS entries",
		Long:  "Print list of reverse DNS entries in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			rdnsList, err := app.client.RDnsGetList()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"ip", "ptr"})

			for _, rdns := range rdnsList {
				t.AppendRow(table.Row{
					rdns.IP,
					rdns.Ptr,
				})
			}

			t.AppendFooter(table.Row{"", "", "Total", len(rdnsList)})
			t.Render()
		},
	}
}

func (app *RobotApp) NewRdnsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rdns:get",
		Short: "Print single reverse DNS entry",
		Long: `Print details of single reverse DNS entry in hetzner account
		reverse DNS entry can be chosen interactively`,
		Run: func(cmd *cobra.Command, args []string) {
			rDnsList, err := app.client.RDnsGetList()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			prompt := promptui.Select{
				Label:             "Select reverse DNS entry",
				Items:             rDnsList,
				Searcher:          getRDnsSearcher(rDnsList),
				Size:              10,
				Templates:         getRDnsSelectTemplates(),
				StartInSearchMode: true,
			}

			choosenIdx, _, err := prompt.Run()
			if err != nil {
				app.logger.Errorln("Prompt failed: ", err)
				return
			}

			choosenRdns := rDnsList[choosenIdx]
			fmt.Println("Chosen reverse DNS entry: ", choosenRdns.IP)

			// directly print info without additional get as that does not deliver more data
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"field", "value"})
			t.AppendRow(table.Row{"ip", choosenRdns.IP})
			t.AppendRow(table.Row{"ptr", choosenRdns.Ptr})

			t.Render()
		},
	}
}

func getRDnsSearcher(rDnsList []models.Rdns) func(string, int) bool {
	return func(input string, index int) bool {
		rDns := rDnsList[index]
		ip := strings.Replace(strings.ToLower(rDns.IP), " ", "", -1)
		ptr := strings.Replace(strings.ToLower(rDns.Ptr), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(ip, input) || strings.Contains(ptr, input)
	}
}

func getRDnsSelectTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "→ {{ .IP | green }} ({{ .Ptr | yellow }})",
		Inactive: "  {{ .IP | cyan }} ({{ .Ptr | red }})",
		Selected: "→ {{ .IP | cyan }}",
		Details: `
	--------- Selected server ----------
	{{ "IP:" | faint }}	          {{ .IP }}
	{{ "PTR record:" | faint }}	  {{ .Ptr }}`,
	}
}
