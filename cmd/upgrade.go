package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
)

var upgradeCommand = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade all commands from remote",
	Run: func(cmd *cobra.Command, args []string) {
		doUpgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCommand)
}

func doUpgrade() {
	if env.Offline {
		fmt.Println("[tips] offline mode can't update from remote.")
		return
	}
	meta, err := fetchLatestCommandMeta()
	if err != nil {
		fmt.Printf("[sorry] failed to fetch latest command metadata, details is: %s\n", err.Error())
		return
	}
	cache.LatestVersion = meta.GetLatestVersion()
	cache.Cmds = meta.GetCommandMaps()
	fetchFileAndFillCache()
	persistCache()
}

func fetchLatestCommandMeta() (*CommandMeta, error) {
	resp, err := http.Get(latestCommandMetaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote responded with %s", resp.Status)
	}
	meta := new(CommandMeta)
	if err := json.NewDecoder(resp.Body).Decode(meta); err != nil {
		return nil, err
	}
	return meta, nil
}

// fetchFileAndFillCache 依次发起 http 请求，将每个 command 对应的 .md 文件缓存到本地
func fetchFileAndFillCache() {
	var all, failed int64
	ch := make(chan *Cmd, 4)
	cmds := cache.GetCmds()
	wg := sync.WaitGroup{}

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Go(func() {
			for item := range ch {
				if err := item.FillSelf(cache.GetLatestVersion()); err != nil {
					atomic.AddInt64(&failed, 1)
				}
				atomic.AddInt64(&all, 1)
				fmt.Printf("[busy working] upgrade command:<%d/%d> => %s\n", atomic.LoadInt64(&all), len(cmds), item.Name)
			}
		})
	}

	for _, item := range cmds {
		ch <- item
	}
	close(ch)
	wg.Wait()
	fmt.Printf("[clap] all commands are upgraded. All: %d, Failed: %d\n", all, failed)
}
