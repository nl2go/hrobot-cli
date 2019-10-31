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
	"gitlab.com/nl2go/hrobot-cli/client/models"
	"gitlab.com/nl2go/hrobot-cli/config"
)

func NewServerGetListCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server:list",
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

			t.AppendHeader(table.Row{"id", "ip", "name", "datacenter", "cancelled"})

			for _, server := range servers {
				t.AppendRow(table.Row{
					server.ServerNumber,
					server.ServerIP,
					server.ServerName,
					server.Dc,
					server.Cancelled,
				})
			}

			t.AppendFooter(table.Row{"", "", "", "Total", len(servers)})
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

			prompt := promptui.Select{
				Label:     "Select server",
				Items:     servers,
				Searcher:  getServerSearcher(servers),
				Size:      10,
				Templates: getServerSelectTemplates(),
			}

			choosenIdx, _, err := prompt.Run()
			if err != nil {
				logger.Errorln("Prompt failed: ", err)
				return
			}

			choosenServer := servers[choosenIdx]
			fmt.Println("Chosen server: ", choosenServer.ServerIP)

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

func NewServerReversalCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server:reverse",
		Short: "Revert single server order",
		Long: `Revert single server order in hetzner account
		server can be chosen interactively`,
		Run: func(cmd *cobra.Command, args []string) {
			robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
			servers, err := robotClient.ServerGetList()
			if err != nil {
				logger.Errorln(err)
				return
			}

			prompt := promptui.Select{
				Label:     "Select server",
				Items:     servers,
				Searcher:  getServerSearcher(servers),
				Size:      10,
				Templates: getServerSelectTemplates(),
			}

			choosenIdx, _, err := prompt.Run()
			if err != nil {
				logger.Errorln("Prompt failed: ", err)
				return
			}

			choosenServer := servers[choosenIdx]
			fmt.Println("Chosen server: ", choosenServer.ServerIP)

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Really reverse server %s (%s) ", choosenServer.ServerName, choosenServer.ServerIP),
				IsConfirm: true,
			}

			confirmRes, err := confirmPrompt.Run()
			if err != nil {
				logger.Errorln("Prompt failed: ", err)
				return
			}

			fmt.Println("Your choice:", confirmRes)

			reverseErr := robotClient.ServerReverse(choosenServer.ServerIP)
			if reverseErr != nil {
				logger.Errorln("Error while reversing server:", reverseErr)
				return
			}

			fmt.Println("Server reversed successfully.")
		},
	}
}

func NewServerSetNamesEmptyCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server:set-names-empty",
		Short: "Sets name for all servers with empty name",
		Long:  "Sets name for all servers with empty name in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
			servers, err := robotClient.ServerGetList()
			if err != nil {
				logger.Errorln(err)
				return
			}

			var serversEmptyName []models.Server
			for _, server := range servers {
				if server.ServerName == "" {
					serversEmptyName = append(serversEmptyName, server)
				}
			}

			if len(serversEmptyName) <= 0 {
				logger.Println("No server with empty name found. Exiting.")
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"id", "ip", "datacenter"})

			for _, server := range serversEmptyName {
				t.AppendRow(table.Row{
					server.ServerNumber,
					server.ServerIP,
					server.Dc,
				})
			}

			t.AppendFooter(table.Row{"", "Total", len(serversEmptyName)})
			t.SetCaption("Servers with empty names")
			t.Render()

			prompt := promptui.Prompt{
				Label: "Add server name prefix",
			}

			prefix, err := prompt.Run()

			if err != nil {
				logger.Errorln("Prompt failed: ", err)
				return
			}

			fmt.Println("Chosen server prefix:", prefix)

			tNames := table.NewWriter()
			tNames.SetOutputMirror(os.Stdout)

			tNames.AppendHeader(table.Row{"id", "ip", "datacenter", "new name"})

			for _, server := range serversEmptyName {
				tNames.AppendRow(table.Row{
					server.ServerNumber,
					server.ServerIP,
					server.Dc,
					generateServerName(server, prefix),
				})
			}

			tNames.AppendFooter(table.Row{"", "", "Total", len(serversEmptyName)})
			tNames.SetCaption("Servers with generated names")
			tNames.Render()

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Really set names as shown above for %d servers", len(serversEmptyName)),
				IsConfirm: true,
			}

			confirmRes, err := confirmPrompt.Run()
			if err != nil {
				logger.Errorln("Prompt failed: ", err)
				return
			}

			fmt.Println("Your choice:", confirmRes)

			for _, server := range serversEmptyName {
				fmt.Println("Set server name for", server.ServerIP, "to", generateServerName(server, prefix), "...")
				err := robotClient.ServerSetName(server.ServerIP, generateServerName(server, prefix))
				if err != nil {
					logger.Errorln(err)
					continue
				}
			}
		},
	}
}

func NewServerGenerateAnsibleInventoryCmd(logger *log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server:gen-ansible-inv",
		Short: "Generates ansible inventory from server list",
		Long:  "Generates ansible inventory from servers in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			robotClient := client.NewBasicAuthClient(cfg.User, cfg.Password)
			servers, err := robotClient.ServerGetList()
			if err != nil {
				logger.Errorln(err)
				return
			}

			invDcs := make(map[string][]string)
			invGroups := make(map[string][]string)

			fmt.Println("[servers]")
			for _, server := range servers {
				dc := strings.Split(server.Dc, "-")
				dcLocation := strings.ToLower(dc[0])

				srvInf := strings.Split(server.ServerName, "-")
				hostGroup := srvInf[0] // first token should define host group (=purpose) of the host, i.e. mongodb

				invDcs[dcLocation] = append(invDcs[dcLocation], server.ServerName)
				invGroups[hostGroup] = append(invGroups[hostGroup], server.ServerName)

				fmt.Printf("%s ansible_host=%s\n", server.ServerName, server.ServerIP)
			}

			for invGroup, groupServers := range invGroups {
				fmt.Printf("\n[%s]\n", invGroup)
				for _, groupServer := range groupServers {
					fmt.Println(groupServer)
				}
			}

			for invDc, dcServers := range invDcs {
				fmt.Printf("\n[dc-%s]\n", invDc)
				for _, dcServer := range dcServers {
					fmt.Println(dcServer)
				}
			}
		},
	}
}

func generateServerName(server models.Server, prefix string) string {
	return fmt.Sprintf("%s-%s-%s-%d", prefix, strings.ToLower(server.Product), strings.ToLower(server.Dc), server.ServerNumber)
}

func getServerSearcher(servers []models.Server) func(string, int) bool {
	return func(input string, index int) bool {
		server := servers[index]
		name := strings.Replace(strings.ToLower(server.ServerName), " ", "", -1)
		ip := strings.Replace(strings.ToLower(server.ServerIP), " ", "", -1)
		product := strings.Replace(strings.ToLower(server.Product), " ", "", -1)
		dc := strings.Replace(strings.ToLower(server.Dc), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input) || strings.Contains(ip, input) || strings.Contains(product, input) || strings.Contains(dc, input)
	}
}

func getServerSelectTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
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
}
