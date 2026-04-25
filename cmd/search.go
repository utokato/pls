package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var searchCommand = &cobra.Command{
	Use:   "search <command>",
	Short: "Search command by keywords",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !cacheReady() {
			printCacheMissingTips()
			return
		}
		doSearch(args[0])
	},
}

func init() {
	rootCmd.AddCommand(searchCommand)
}

func doSearch(keyword string) {
	table := tablewriter.NewTable(
		os.Stdout,
		tablewriter.WithRendition(tw.Rendition{
			Settings: tw.Settings{
				Separators: tw.Separators{BetweenRows: tw.On},
			},
		}),
	)
	table.Header([]string{"command", "description"})
	keyword = strings.ToLower(keyword)
	for k, v := range cache.GetCmds() {
		k = strings.ToLower(k)
		if strings.Contains(k, keyword) {
			table.Append([]string{v.Name, v.Desc})
			continue
		}
		desc := strings.ToLower(v.Desc)
		if strings.Contains(desc, keyword) {
			table.Append([]string{v.Name, v.Desc})
			continue
		}
	}
	if err := table.Render(); err != nil {
		fmt.Println("[sorry] failed to render search table")
	}
}
