package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"pls/offline"
)

var offlineCommand = &cobra.Command{
	Use:   "offline <enable|disable>",
	Short: "Set the current context to offline, and load infos from local",
	Example: `# Set offline ctx
	pls offline enable
					
# Set online ctx
	pls offline disable`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flag, ok := parseOfflineFlag(args[0])
		if !ok {
			fmt.Println("[sorry] offline command only accepts enable or disable")
			return
		}
		doSetOffline(flag)
	},
}

func init() {
	rootCmd.AddCommand(offlineCommand)
}

func parseOfflineFlag(raw string) (bool, bool) {
	switch raw {
	case "enable":
		return true, true
	case "disable":
		return false, true
	default:
		return false, false
	}
}

func doSetOffline(flag bool) {
	if !flag {
		// 如果关闭离线模式，则先清理之前的旧数据
		err := os.RemoveAll(dirPath)
		if err != nil {
			fmt.Println("[sorry] clear data fail when set offline closed")
			return
		}
		err = makeCmdDir(dirPath)
		if err != nil {
			fmt.Println("[sorry] failed to recreate command dir when set offline closed")
			return
		}
		persistEnv(false, false)
		return
	}
	// 开启离线模式，但是之前的压缩包已经解压过，就不在处理了
	if env.Decompressed {
		persistEnv(true, true)
		return
	}

	readerAt := bytes.NewReader(offline.Resource)
	// 将离线资源的压缩包，解压缩到指定目录
	archive, err := zip.NewReader(readerAt, int64(readerAt.Len()))
	if err != nil {
		fmt.Printf("[sorry] unzip offline resources met error, details is: %s\n", err.Error())
		return
	}

	// 构建 cache
	commandMaps := make(map[string]*Cmd, 1024)
	for _, f := range archive.File {
		name := removeExt(f.Name)
		cmd := &Cmd{
			Name: name,
		}
		reader, err := f.Open()
		if err != nil {
			fmt.Printf("[sorry] open offline file met error, details is: %s\n", err.Error())
			return
		}
		// 持久化 静态 md 资源文件
		err = cmd.FillSelfByReader(reader)
		if err != nil {
			fmt.Printf("[sorry] write offline file met error, details is: %s\n", err.Error())
			return
		}
		fmt.Printf("[busy working] release command info => %s\n", name)
		commandMaps[name] = cmd
	}

	cache.LatestVersion = version
	cache.Cmds = commandMaps
	// 持久化 cache
	persistCache()

	// 持久化 env
	persistEnv(flag, true)
}

// removeExt 去除文件名的扩展名
func removeExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
