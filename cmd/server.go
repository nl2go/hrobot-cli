package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.com/nl2go/hrobot-cli/client"
	"gitlab.com/nl2go/hrobot-cli/config"
)

func NewServerGetListCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
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

			for _, server := range servers {
				t.AppendRow(table.Row{
					server.ServerNumber,
					server.ServerIP,
					server.ServerName,
					server.Dc,
				})
			}

			t.AppendFooter(table.Row{"", "", "Total", len(servers)})
			t.Render()
		},
	}
}

func NewServerGetCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server:get",
		Short: "Print single server",
		Long: `Print details of single server in hetzner account
		server can be chosen interactively`,
		Run: func(cmd *cobra.Command, args []string) {
			robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
			servers, err := robotClient.ServerGetList()
			if err != nil {
				logger.Errorln(err)
				return
			}

			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "→ {{ .ServerName | cyan }} ({{ .ServerIP | red }} - {{ .Product | blue }} - {{ .Dc | yellow }})",
				Inactive: "  {{ .ServerName | cyan }} ({{ .ServerIP | red }} - {{ .Product | blue }} - {{ .Dc | yellow }})",
				Selected: "→ {{ .ServerName | red | cyan }}",
				Details: `
			--------- Selected server ----------
			{{ "Server name:" | faint }}  {{ .ServerName }}
			{{ "IP:" | faint }}	          {{ .ServerIP }}
			{{ "Product:" | faint }}	  {{ .Product }}
			{{ "Datacenter:" | faint }}	  {{ .Dc }}`,
			}

			searcher := func(input string, index int) bool {
				server := servers[index]
				name := strings.Replace(strings.ToLower(server.ServerName), " ", "", -1)
				ip := strings.Replace(strings.ToLower(server.ServerIP), " ", "", -1)
				product := strings.Replace(strings.ToLower(server.Product), " ", "", -1)
				dc := strings.Replace(strings.ToLower(server.Dc), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input) || strings.Contains(ip, input) || strings.Contains(product, input) || strings.Contains(dc, input)
			}

			prompt := promptui.Select{
				Label:     "Select server",
				Items:     servers,
				Searcher:  searcher,
				Size:      10,
				Templates: templates,
			}

			choosenIdx, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			choosenServer := servers[choosenIdx]
			fmt.Println("You choose: ", choosenServer.ServerIP)

			server, err := robotClient.ServerGet(choosenServer.ServerIP)
			if err != nil {
				logger.Errorln(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"field", "value"})
			t.AppendRow(table.Row{"number", server.ServerNumber})
			t.AppendRow(table.Row{"ip", server.ServerIP})
			t.AppendRow(table.Row{"data center", server.Dc})
			t.AppendRow(table.Row{"product", server.Product})
			t.AppendRow(table.Row{"status", server.Status})
			t.AppendRow(table.Row{"subnet", "IP: " + server.Subnet[0].IP + " Mask: " + server.Subnet[0].Mask})
			t.AppendRow(table.Row{"traffic", server.Traffic})
			t.AppendRow(table.Row{"paid until", server.PaidUntil})

			t.Render()
		},
	}
}
