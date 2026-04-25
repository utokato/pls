package cmd

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/spf13/cobra"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path"
	"pls/offline"
	"pls/resp"
	"strings"
	"syscall"
)

const port = 6023

var serveCommand = &cobra.Command{
	Use:    "serve",
	Short:  "Serve a web server",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		doServe()
	},
}

type SearchItem struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func init() {
	rootCmd.AddCommand(serveCommand)
}

func doServe() {
	server := echo.New()
	server.Use(middleware.Recover())
	server.Use(middleware.CORS("*"))

	server.GET("/v1/healthz", handleHealthz)
	server.GET("/v1/command/search", handleSearch)
	server.GET("/v1/command/show", handleShow)

	// Serve frontend resources
	handleStaticFilesServer(server)

	// Start server
	address := fmt.Sprintf(":%d", port)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := (echo.StartConfig{Address: address, HideBanner: true}).Start(ctx, server); err != nil {
		fmt.Println("[sorry] failed to start server")
	}
}

func handleStaticFilesServer(e *echo.Echo) {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: getFileSystem("dist"),
	}))
	e.StaticFS(
		"/assets",
		getFileSystem("dist/assets"),
		middleware.GzipWithConfig(
			middleware.GzipConfig{
				Level: 5,
			}),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c *echo.Context) error {
				c.Response().Header().Set(echo.HeaderCacheControl, "max-age=31536000, immutable")
				return next(c)
			}
		},
	)
}

func handleShow(c *echo.Context) error {
	cmdName := c.QueryParam("keyword")
	cmds := cache.GetCmds()
	command, exist := cmds[cmdName]
	if !exist {
		msg := fmt.Sprintf("[sorry] could not found command <%s>", cmdName)
		return c.JSON(http.StatusNotFound, resp.Failure(msg))
	}
	// 将 .md 文件内容返回给界面
	fp := path.Join(dirPath, fmt.Sprintf("%s.md", command.Name))
	if !fileExist(fp) {
		msg := fmt.Sprintf("[sorry] could not found command <%s>", cmdName)
		return c.JSON(http.StatusNotFound, resp.Failure(msg))
	}
	source, err := os.ReadFile(fp)
	if err != nil {
		msg := fmt.Sprintf("[sorry] failed to open file <%s>\n", fp)
		return c.JSON(http.StatusNotFound, resp.Failure(msg))
	}

	return c.JSON(http.StatusOK, resp.Success(string(source)))
}

func handleSearch(c *echo.Context) error {
	var res []SearchItem
	key := c.QueryParam("keyword")
	if key == "" {
		return c.JSON(http.StatusOK, resp.Success(res))
	}
	// search
	key = strings.ToLower(key)
	for k, v := range cache.GetCmds() {
		k = strings.ToLower(k)
		if strings.Contains(k, key) {
			res = append(res, SearchItem{
				Name: v.Name,
				Desc: v.Desc,
			})
			continue
		}
		desc := strings.ToLower(v.Desc)
		if strings.Contains(desc, key) {
			res = append(res, SearchItem{
				Name: v.Name,
				Desc: v.Desc,
			})
			continue
		}
	}
	return c.JSON(http.StatusOK, resp.Success(res))
}

func handleHealthz(c *echo.Context) error {
	return c.JSON(http.StatusOK, resp.Success("Ok"))
}

func getFileSystem(path string) fs.FS {
	f, err := fs.Sub(offline.Dist, path)
	if err != nil {
		panic(err)
	}
	return f
}
