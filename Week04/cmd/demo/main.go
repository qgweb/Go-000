package main

import (
	"Go-000/Week04/internal/conf"
	xhttp "Go-000/Week04/internal/server/http"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configPath = flag.String("conf", "./config/conf.toml", "配置文件路径")

	httpServer *http.Server
	webServer  *xhttp.WebServer
)

func init() {
	flag.Parse()
}

func main() {
	config, err := conf.Decode(*configPath)
	if err != nil {
		panic(err)
	}

	httpServer = &http.Server{Addr: config.Base.HttpAddr, Handler: nil}
	webServer = InitializeServer(config, httpServer)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("开始监听：", config.Base.HttpAddr)
		if err := httpServer.ListenAndServe(); err != nil {
			cancel()
			log.Println("监听失败:", err)
		}
	}()
	handleSignal(ctx)
}

func handleSignal(ctx context.Context) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sc:
			fmt.Println(sig, "收到信号")
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				fmt.Println("开始退出进程")
				ctx1, _ := context.WithTimeout(context.Background(), time.Second*5)
				webServer.Close(ctx1)
				os.Exit(0)
			case syscall.SIGHUP:
				fmt.Println("开始平滑重启")
			default:
				fmt.Println("未知信号量")
			}
		case <-ctx.Done():
			fmt.Println("退出监听")
		}
	}
}
