package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/nl2go/hrobot-go/models"
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
			chosenServer, err := app.selectServer()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			// additional get as getting a single server returns more data
			server, err := app.client.ServerGet(chosenServer.ServerIP)
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
			chosenServer, err := app.selectServer()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Really reverse server %s (%s) ", chosenServer.ServerName, chosenServer.ServerIP),
				IsConfirm: true,
			}

			_, confirmErr := confirmPrompt.Run()
			if confirmErr != nil {
				app.logger.Errorln("Prompt failed: ", confirmErr)
				return
			}

			_, reverseErr := app.client.ServerReverse(chosenServer.ServerIP)
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
			chosenServers, err := app.selectMultipleServers()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			color.Cyan("Servers selected for renaming:")

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"name", "ip"})

			for _, chServ := range chosenServers {
				t.AppendRow(table.Row{
					chServ.ServerName,
					chServ.ServerIP,
				})
			}

			t.AppendFooter(table.Row{"Total", len(chosenServers)})
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

			for _, server := range chosenServers {
				tNames.AppendRow(table.Row{
					server.ServerNumber,
					server.ServerIP,
					server.ServerName,
					generateServerName(server, prefix),
				})
			}

			tNames.AppendFooter(table.Row{"", "", "Total", len(chosenServers)})
			tNames.SetCaption("Servers with generated names")
			tNames.Render()

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Really set names as shown above for %d servers", len(chosenServers)),
				IsConfirm: true,
			}

			_, confirmErr := confirmPrompt.Run()
			if confirmErr != nil {
				app.logger.Errorln("Prompt failed: ", confirmErr)
				return
			}

			for _, server := range chosenServers {
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

func (app *RobotApp) NewServerActivateRescueCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server:rescue",
		Short: "Activate rescue mode for single server",
		Long:  `Activate rescue mode for single server in hetzner account, server can be chosen interactively`,
		Run: func(cmd *cobra.Command, args []string) {
			chosenServer, err := app.selectServer()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			rescueOptions, rescueOptErr := app.client.BootRescueGet(chosenServer.ServerIP)
			if rescueOptErr != nil {
				app.logger.Errorln("Error while fetching rescue options:", rescueOptErr)
				return
			}

			selectOsItems := rescueOptions.Os.([]interface{})

			promptRescueOS := promptui.Select{
				Label: "Select rescue operating system",
				Items: selectOsItems,
				Size:  10,
			}

			// convert interface to string
			selectOs := make([]string, len(selectOsItems))
			for i, v := range selectOsItems {
				selectOs[i] = v.(string)
			}

			chosenOSIdx, _, err := promptRescueOS.Run()
			if err != nil {
				app.logger.Errorln("Prompt failed: ", err)
				return
			}

			chosenOS := selectOs[chosenOSIdx]
			color.Cyan(fmt.Sprint("Chosen OS: ", chosenOS))

			selectArchItems := rescueOptions.Arch.([]interface{})

			promptRescueArch := promptui.Select{
				Label: "Select rescue operating system architecture",
				Items: selectArchItems,
			}

			// convert interface to int
			selectArch := make([]int, len(selectArchItems))
			for i, v := range selectArchItems {
				selectArch[i] = int(v.(float64))
			}

			chosenArchIdx, _, err := promptRescueArch.Run()
			if err != nil {
				app.logger.Errorln("Prompt failed: ", err)
				return
			}

			chosenArch := selectArch[chosenArchIdx]
			color.Cyan(fmt.Sprint("Chosen arch: ", chosenArch))

			confirmPromptKey := promptui.Prompt{
				Label:     fmt.Sprintf("Use SSH key for rescue system "),
				IsConfirm: true,
				Default:   "y",
			}

			useSSHKey := true
			_, confirmKeyErr := confirmPromptKey.Run()
			if confirmKeyErr != nil {
				useSSHKey = false
			}

			var chosenKey models.Key
			if useSSHKey {
				keys, err := app.client.KeyGetList()
				if err != nil {
					app.logger.Errorln(err)
					return
				}

				promptKey := promptui.Select{
					Label:     "Select key",
					Items:     keys,
					Size:      10,
					Templates: getKeySelectTemplates(),
				}

				chosenKeyIdx, _, err := promptKey.Run()
				if err != nil {
					app.logger.Errorln("Prompt failed: ", err)
					return
				}

				chosenKey = keys[chosenKeyIdx]
				color.Cyan(fmt.Sprint("Chosen key: ", chosenKey.Name, chosenKey.Fingerprint))
			} else {
				color.Cyan("Chosen to use password instead of key.")
			}

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Really activate rescue system and reboot server %s (%s) ", chosenServer.ServerName, chosenServer.ServerIP),
				IsConfirm: true,
			}

			_, confirmErr := confirmPrompt.Run()
			if confirmErr != nil {
				app.logger.Errorln("Prompt failed: ", confirmErr)
				return
			}

			input := &models.RescueSetInput{}
			input = &models.RescueSetInput{
				OS:            chosenOS,
				Arch:          chosenArch,
				AuthorizedKey: chosenKey.Fingerprint,
			}

			rescue, err := app.client.BootRescueSet(chosenServer.ServerIP, input)
			if err != nil {
				app.logger.Errorln("Error while activating rescue system:", err)
				return
			}

			resetInput := &models.ResetSetInput{
				Type: models.ResetTypeHardware,
			}

			_, resetErr := app.client.ResetSet(chosenServer.ServerIP, resetInput)
			if resetErr != nil {
				app.logger.Errorln("Error while rebooting server:", resetErr)
				return
			}

			if !useSSHKey {
				color.Cyan(fmt.Sprintf("Password for accessing rescue mode: %s", rescue.Password))
			}

			color.Cyan("Rescue mode successfully activated and server rebooted.")
		},
	}
}

func (app *RobotApp) NewServerResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server:reset",
		Short: "Reset single server (hardware reset)",
		Long:  `Reset single server in hetzner account using hardware reset, server can be chosen interactively`,
		Run: func(cmd *cobra.Command, args []string) {
			chosenServer, err := app.selectServer()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Really reset server %s (%s) ", chosenServer.ServerName, chosenServer.ServerIP),
				IsConfirm: true,
			}

			_, confirmErr := confirmPrompt.Run()
			if confirmErr != nil {
				app.logger.Errorln("Prompt failed: ", confirmErr)
				return
			}

			resetInput := &models.ResetSetInput{
				Type: models.ResetTypeHardware,
			}

			_, resetErr := app.client.ResetSet(chosenServer.ServerIP, resetInput)
			if resetErr != nil {
				app.logger.Errorln("Error while rebooting server:", resetErr)
				return
			}

			color.Cyan("Server rebooted successfully.")
		},
	}
}

func (app *RobotApp) NewServerGenerateAnsibleInventoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server:ansible-inv",
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

func (app *RobotApp) selectServer() (*models.Server, error) {
	servers, err := app.client.ServerGetList()
	if err != nil {
		return nil, err
	}

	prompt := promptui.Select{
		Label:             "Select server",
		Items:             servers,
		Searcher:          getServerSearcher(servers),
		Size:              10,
		Templates:         getServerSelectTemplates(),
		StartInSearchMode: true,
	}

	chosenIdx, _, err := prompt.Run()
	if err != nil {
		app.logger.Errorln("Prompt failed: ", err)
		return nil, err
	}

	chosenServer := servers[chosenIdx]
	color.Cyan(fmt.Sprint("Chosen server: ", chosenServer.ServerIP))

	return &chosenServer, nil
}

func (app *RobotApp) selectMultipleServers() ([]models.Server, error) {
	var selectServers []models.Server

	selectServers = append(selectServers, models.Server{
		ServerName: "Done",
	})

	servers, err := app.client.ServerGetList()
	if err != nil {
		return []models.Server{}, err
	}

	for _, server := range servers {
		selectServers = append(selectServers, server)
	}

	chosenIdx := 1
	var chosenServers []models.Server

	for chosenIdx > 0 {
		if len(chosenServers) > 0 {
			color.Cyan("Servers  currently selected:")

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"name", "ip"})

			for _, chServ := range chosenServers {
				t.AppendRow(table.Row{
					chServ.ServerName,
					chServ.ServerIP,
				})
			}

			t.AppendFooter(table.Row{"Total", len(chosenServers)})
			t.Render()
		} else {
			color.Cyan("No servers currently selected.")
		}

		prompt := promptui.Select{
			Label:             "Choose servers",
			Items:             selectServers,
			Searcher:          getServerSearcher(selectServers),
			Size:              10,
			Templates:         getServerSelectTemplates(),
			StartInSearchMode: true,
		}

		chosenIdx, _, err = prompt.Run()
		if err != nil {
			app.logger.Errorln("Prompt failed: ", err)
			return []models.Server{}, err
		}

		if chosenIdx > 0 {
			chosenServer := selectServers[chosenIdx]
			color.Cyan(fmt.Sprint("Chosen server: ", chosenServer.ServerIP))

			selectedBefore := false
			for _, chServ := range chosenServers {
				if chServ.ServerIP == chosenServer.ServerIP {
					selectedBefore = true
				}
			}

			if selectedBefore {
				app.logger.Infof("Server %s was selected already.", chosenServer.ServerIP)
			} else {
				chosenServers = append(chosenServers, chosenServer)
				// remove chosen server from select list = re-slicing
				selectServers = append(selectServers[:chosenIdx], selectServers[chosenIdx+1:]...)
			}
		}
	}

	return chosenServers, nil
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
		Label:    "{{ . }} ?",
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
