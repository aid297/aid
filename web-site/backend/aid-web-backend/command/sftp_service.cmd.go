package command

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
)

type SFTPServiceCommand struct{}

func (*SFTPServiceCommand) Launch() {
	var (
		err    error
		port   = flag.String("port", "8080", "监听端口，如8080、9000") // 定义命令行参数：端口（默认8080）、共享目录（默认当前目录）
		absDir string
	)
	flag.Parse()

	if port == nil || *port == "" {
		port = &global.CONFIG.FileManager.Port
	}

	// 验证共享目录是否存在
	if absDir, err = filepath.Abs(global.CONFIG.FileManager.Dir); err != nil {
		global.LOG.Error("目录路径错误", zap.Error(err))
		return
	}
	if _, err = os.Stat(absDir); os.IsNotExist(err) {
		global.LOG.Error("目录不存在", zap.String("dir", absDir), zap.Error(err))
		return
	}

	// 获取Mac的局域网方便虚拟桌面访问
	localIP := getLocalIP()
	if localIP == "" {
		global.LOG.Warn("⚠️ 未检测到局域网IP，请手动确认Mac的IP地址")
	} else {
		global.LOG.Info("文件服务器已启动", zap.String("共享目录", global.CONFIG.FileManager.Dir))
		fmt.Printf("✅ 文件服务器已启动：\n")
		fmt.Printf("   共享目录：%s\n", absDir)
		fmt.Printf("   访问地址：http://%s:%s\n", localIP, *port)
		fmt.Printf("   本地访问：http://localhost:%s\n", *port)
		fmt.Println("📌 提示：保持终端窗口打开，关闭则停止服务")
	}

	// 启动HTTP文件服务器，支持目录浏览和文件下载
	http.Handle("/", http.FileServer(http.Dir(absDir)))

	// 监听指定端口（0.0.0.0表示允许所有IP访问）
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", *port), nil))
}

// 获取Mac的局域网IP（排除回环地址127.0.0.1）
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
