# pls — Linux 命令速查 CLI 工具

> 基于 https://github.com/chenjiandongx/pls 二次开发，新增离线模式与内嵌 Web 界面。

## 项目结构

```
pls/
├── main.go          # 程序入口
├── build.sh         # 一键构建脚本（前端 + Go）
├── go.mod
├── cmd/             # CLI 命令实现（Cobra）
│   ├── root.go      # 根命令、初始化、工具函数
│   ├── search.go    # search 命令
│   ├── show.go      # show 命令
│   ├── serve.go     # serve 命令（内嵌 Web 服务器）
│   ├── upgrade.go   # upgrade 命令
│   ├── offline.go   # offline 命令
│   ├── clear.go     # clear 命令
│   ├── version.go   # version 命令
│   ├── model.go     # 数据模型与 HTTP 工具
│   └── tool.go      # 本地文件/持久化工具函数
├── offline/
│   ├── resource.go  # 通过 embed 内嵌 offline.zip 和前端 dist
│   └── offline.zip  # 离线命令资源包
├── resp/
│   └── resp.go      # HTTP API 统一响应结构
├── web/             # Vue 3 前端（Vite 构建）
│   ├── src/
│   └── ...
└── bin/
    ├── pls          # 构建产物（Linux amd64）
    └── start.sh     # 安装并启动脚本
```

## 构建

需要提前安装 [pnpm](https://pnpm.io/)。

```shell
# 在 Windows 上交叉编译为 Linux 二进制
./build.sh
```

`build.sh` 会依次执行：
1. `pnpm build` 编译前端，输出到 `offline/dist/`
2. 交叉编译 Go 程序，输出到 `bin/pls`

## 安装

```shell
# 将二进制安装到系统并以后台模式启动 Web 服务
./bin/start.sh
```

或手动安装：

```shell
cp bin/pls /usr/local/bin/pls
chmod +x /usr/local/bin/pls
```

## 命令用法

### 初始化资源

```shell
# 在线模式：从 unpkg.com 拉取最新命令数据
pls upgrade

# 离线模式：从内嵌资源包解压
pls offline enable

# 切换回在线模式（会清除本地缓存）
pls offline disable
```

### 搜索与查看

```shell
# 按关键字搜索命令（支持中英文）
pls search 压缩
pls search zip

# 查看命令详情（Markdown 渲染到终端）
pls show zip

# 强制从远程刷新后查看（仅在线模式）
pls show zip -f
```

### Web 界面

```shell
pls serve
```

访问 http://localhost:6023

后台运行：

```shell
nohup pls serve >> pls.log 2>&1 &
```

### 其他

```shell
# 查看版本
pls version

# 清除所有本地缓存
pls clear
```

