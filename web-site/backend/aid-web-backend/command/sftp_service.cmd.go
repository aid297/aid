package command

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/aid297/aid/debugLogger"
	"github.com/aid297/aid/filesystem/filesystemV4"
	"github.com/aid297/aid/str"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
)

type SFTPServiceCommand struct{}

func (*SFTPServiceCommand) Launch() {
	type output struct {
		Dir  string
		IP   string
		Port string
	}

	var (
		port       = flag.String("port", "8080", "监听端口，如8080、9000") // 定义命令行参数：端口（默认8080）、共享目录（默认当前目录）
		outputTemp *str.Template[output]
		dir        filesystemV4.IFilesystem
	)
	flag.Parse()

	if port == nil || *port == "" {
		port = &global.CONFIG.FileManager.Port
	}

	if dir = filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.FileManager.Dir)); !dir.GetExist() {
		dir.Create()
	}

	// 获取Mac的局域网方便虚拟桌面访问
	localIP := getLocalIP()
	if localIP == "" {
		global.LOG.Warn("⚠️ 未检测到局域网IP，请手动确认 Mac 的 IP 地址")
		return
	}

	if outputTemp = str.NewTemplate("文件服务器已启动", `
✅  文件服务器已启动：
    共享目录：{{.Dir}}
    访问地址：http://{{.IP}}:{{.Port}}
    本地访问：http://localhost:{{.Port}}
📌  提示：保持终端窗口打开，关闭则停止服务
		`, output{Dir: dir.GetFullPath(), IP: localIP, Port: *port}); outputTemp.Error() != nil {
		global.LOG.Error("生成输出字符串失败", zap.Error(outputTemp.Error()))
		return
	}
	debugLogger.Print(outputTemp.String())

	// 启动HTTP文件服务器，支持目录浏览和文件下载
	http.Handle("/", http.FileServer(http.Dir(dir.GetFullPath())))

	// 监听指定端口（0.0.0.0 表示允许所有 IP 访问）
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", *port), nil))
}

// 获取 Mac 的局域网IP（排除回环地址 127.0.0.1）
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		if ipNet.IP.To4() != nil { // 只返回IPv4地址（虚拟桌面更兼容）
			return ipNet.IP.String()
		}
	}
	return ""
}
