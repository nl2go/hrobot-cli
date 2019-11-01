package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"gitlab.com/newsletter2go/hrobot-go/models"
)

func (app *RobotApp) NewServerGetListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server:list",
		Short: "Print list of servers",
		Long:  "Print list of servers in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			servers, err := app.client.ServerGetList()
			if err != nil {
				app.logger.Errorln(err)
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

func (app *RobotApp) NewServerGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server:get",
		Short: "Print single server",
		Long: `Print details of single server in hetzner account
		server can be chosen interactively`,
		Run: func(cmd *cobra.Command, args []string) {
			servers, err := app.client.ServerGetList()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			prompt := promptui.Select{
				Label:             "Select server",
				Items:             servers,
				Searcher:          getServerSearcher(servers),
				Size:              10,
				Templates:         getServerSelectTemplates(),
				StartInSearchMode: true,
			}

			choosenIdx, _, err := prompt.Run()
			if err != nil {
				app.logger.Errorln("Prompt failed: ", err)
				return
			}

			choosenServer := servers[choosenIdx]
			color.Cyan(fmt.Sprint("Chosen server: ", choosenServer.ServerIP))

			// additional get as getting a single server returns more data
			server, err := app.client.ServerGet(choosenServer.ServerIP)
			if err != nil {
				app.logger.Errorln(err)
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

func (app *RobotApp) NewServerReversalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server:reverse",
		Short: "Revert single server order",
		Long:  `Revert single server order in hetzner account, server can be chosen interactively`,
		Run: func(cmd *cobra.Command, args []string) {
			servers, err := app.client.ServerGetList()
			if err != nil {
				app.logger.Errorln(err)
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
				app.logger.Errorln("Prompt failed: ", err)
				return
			}

			choosenServer := servers[choosenIdx]
			color.Cyan(fmt.Sprint("Chosen server: ", choosenServer.ServerIP))

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Really reverse server %s (%s) ", choosenServer.ServerName, choosenServer.ServerIP),
				IsConfirm: true,
			}

			_, confirmErr := confirmPrompt.Run()
			if confirmErr != nil {
				app.logger.Errorln("Prompt failed: ", confirmErr)
				return
			}

			_, reverseErr := app.client.ServerReverse(choosenServer.ServerIP)
			if reverseErr != nil {
				app.logger.Errorln("Error while reversing server:", reverseErr)
				return
			}

			color.Cyan("Server reversed successfully.")
		},
	}
}

func (app *RobotApp) NewServerSetNameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server:set-name",
		Short: "Sets name for selected servers",
		Long:  "Sets name for selected servers in the hetzner account, servers can be chosen interactively",
		Run: func(cmd *cobra.Command, args []string) {
			var selectServers []models.Server

			selectServers = append(selectServers, models.Server{
				ServerName: "Done",
			})

			servers, err := app.client.ServerGetList()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			for _, server := range servers {
				selectServers = append(selectServers, server)
			}

			choosenIdx := 1
			var choosenServers []models.Server

			for choosenIdx > 0 {
				if len(choosenServers) > 0 {
					color.Cyan("Servers  currently selected for renaming:")

					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)

					t.AppendHeader(table.Row{"name", "ip"})

					for _, chServ := range choosenServers {
						t.AppendRow(table.Row{
							chServ.ServerName,
							chServ.ServerIP,
						})
					}

					t.AppendFooter(table.Row{"Total", len(choosenServers)})
					t.Render()
				} else {
					color.Cyan("No servers currently selected for renaming.")
				}

				prompt := promptui.Select{
					Label:             "Choose servers for renaming",
					Items:             selectServers,
					Searcher:          getServerSearcher(selectServers),
					Size:              10,
					Templates:         getServerSelectTemplates(),
					StartInSearchMode: true,
				}

				choosenIdx, _, err = prompt.Run()
				if err != nil {
					app.logger.Errorln("Prompt failed: ", err)
					return
				}

				if choosenIdx > 0 {
					choosenServer := selectServers[choosenIdx]
					color.Cyan(fmt.Sprint("Chosen server: ", choosenServer.ServerIP))

					selectedBefore := false
					for _, chServ := range choosenServers {
						if chServ.ServerIP == choosenServer.ServerIP {
							selectedBefore = true
						}
					}

					if selectedBefore {
						app.logger.Infof("Server %s was selected already.", choosenServer.ServerIP)
					} else {
						choosenServers = append(choosenServers, choosenServer)
						// remove chosen server from select list = re-slicing
						selectServers = append(selectServers[:choosenIdx], selectServers[choosenIdx+1:]...)
					}
				}
			}

			color.Cyan("Servers selected for renaming:")

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"name", "ip"})

			for _, chServ := range choosenServers {
				t.AppendRow(table.Row{
					chServ.ServerName,
					chServ.ServerIP,
				})
			}

			t.AppendFooter(table.Row{"Total", len(choosenServers)})
			t.Render()

			prompt := promptui.Prompt{
				Label: "Add server name prefix",
			}

			prefix, err := prompt.Run()

			if err != nil {
				app.logger.Errorln("Prompt failed: ", err)
				return
			}

			color.Cyan(fmt.Sprint("Chosen server prefix: ", prefix))

			tNames := table.NewWriter()
			tNames.SetOutputMirror(os.Stdout)

			tNames.AppendHeader(table.Row{"id", "ip", "current name", "new name"})

			for _, server := range choosenServers {
				tNames.AppendRow(table.Row{
					server.ServerNumber,
					server.ServerIP,
					server.ServerName,
					generateServerName(server, prefix),
				})
			}

			tNames.AppendFooter(table.Row{"", "", "Total", len(choosenServers)})
			tNames.SetCaption("Servers with generated names")
			tNames.Render()

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Really set names as shown above for %d servers", len(choosenServers)),
				IsConfirm: true,
			}

			_, confirmErr := confirmPrompt.Run()
			if confirmErr != nil {
				app.logger.Errorln("Prompt failed: ", confirmErr)
				return
			}

			for _, server := range choosenServers {
				color.Cyan(fmt.Sprint("Set server name for ", server.ServerIP, " to ", generateServerName(server, prefix), " ..."))

				input := &models.ServerSetNameInput{
					Name: generateServerName(server, prefix),
				}

				_, err := app.client.ServerSetName(server.ServerIP, input)
				if err != nil {
					app.logger.Errorln(err)
					continue
				}
			}
		},
	}
}

func (app *RobotApp) NewServerGenerateAnsibleInventoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server:gen-ansible-inv",
		Short: "Generates ansible inventory from server list",
		Long:  "Generates ansible inventory from servers in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			servers, err := app.client.ServerGetList()
			if err != nil {
				app.logger.Errorln(err)
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
	return fmt.Sprintf("%s-%s-hetzner-%s-%d", prefix, strings.ToLower(server.Product), strings.ToLower(server.Dc), server.ServerNumber)
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
		Active:   "→ {{ .ServerName | green }} ({{ .ServerIP | yellow }} - {{ .Product | yellow }} - {{ .Dc | yellow }})",
		Inactive: "  {{ .ServerName | cyan }} ({{ .ServerIP | red }} - {{ .Product | blue }} - {{ .Dc | green }})",
		Selected: "→ {{ .ServerName | cyan }}",
		Details: `
	--------- Selected server ----------
	{{ "Server name:" | faint }}  {{ .ServerName }}
	{{ "IP:" | faint }}	          {{ .ServerIP }}
	{{ "Product:" | faint }}	  {{ .Product }}
	{{ "Datacenter:" | faint }}	  {{ .Dc }}`,
	}
}
