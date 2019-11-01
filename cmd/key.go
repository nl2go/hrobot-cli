package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func (app *RobotApp) NewKeyGetListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "key:list",
		Short: "Print list of ssh keys",
		Long:  "Print list of ssh keys in the hetzner account",
		Run: func(cmd *cobra.Command, args []string) {
			keys, err := app.client.KeyGetList()
			if err != nil {
				app.logger.Errorln(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"name", "type", "size", "fingerprint"})

			for _, key := range keys {
				t.AppendRow(table.Row{
					key.Name,
					key.Type,
					key.Size,
					key.Fingerprint,
				})
			}

			t.AppendFooter(table.Row{"", "", "Total", len(keys)})
			t.Render()
		},
	}
}

func getKeySelectTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }} ?",
		Active:   "→ {{ .Name | green }} ({{ .Fingerprint | yellow }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Fingerprint | red }})",
		Selected: "→ {{ .Name | cyan }}",
	}
}
