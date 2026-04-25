package cmd

import (
	"errors"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var showCommand = &cobra.Command{
	Use:   "show <command>",
	Short: "Show the specified command usage",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !cacheReady() {
			printCacheMissingTips()
			return
		}
		force, _ := cmd.Flags().GetBool("force")
		doShow(args[0], force)
	},
}

func init() {
	rootCmd.AddCommand(showCommand)

	showCommand.Flags().BoolP("force", "f", false, "force to refresh command usage from remote.")
}

func doShow(cmdName string, force bool) {
	cmds := cache.GetCmds()
	command, exist := cmds[cmdName]
	if !exist {
		fmt.Printf("[sorry] could not found command <%s>\n", cmdName)
		return
	}
	// only online mode can force refresh
	if force && !env.Offline {
		if err := command.FillSelf(cache.GetLatestVersion()); err != nil {
			if errors.Is(err, ErrCommandNotFound) {
				fmt.Printf("[sorry] could not found command <%s>\n", cmdName)
			} else {
				fmt.Printf("[sorry] failed to download command <%s>\n", cmdName)
			}
			return
		}
		// 强制更新获取了最新的 command 对应的 .md 文件后，需要重新构建一下本地缓存
		persistCache()
	}
	// 将 .md 文件展示到控制台上
	fp := path.Join(dirPath, fmt.Sprintf("%s.md", command.Name))
	if !fileExist(fp) {
		fmt.Printf("[sorry] could not found command <%s>\n", cmdName)
		return
	}
	source, err := os.ReadFile(fp)
	if err != nil {
		fmt.Printf("[sorry] failed to open file <%s>\n", fp)
		return
	}
	markdown.BlueBgItalic = color.New(color.FgBlue).SprintFunc()
	result := markdown.Render(string(source), 80, 6)
	fmt.Println(string(result))
}
