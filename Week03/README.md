# 作业
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

# 代码
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Second * 10)
		writer.WriteHeader(200)
		_, _ = writer.Write([]byte("ok"))
	})

	server := &http.Server{Addr: ":8080", Handler: mux}
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		fmt.Println("开始监听端口", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			return errors.Wrap(err, "listen server")
		}
		return nil
	})

	g.Go(func() error {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sc:
			fmt.Println(sig, "收到信号")
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				fmt.Println("开始退出进程")
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				if err := server.Shutdown(ctx); err != nil {
					return errors.Wrap(err, "shutdown server")
				} else {
					return nil
				}
			case syscall.SIGHUP:
				fmt.Println("开始平滑重启")
			default:
				fmt.Println("未知信号量")
			}
		case <-ctx.Done():
			fmt.Println("退出监听")
			return nil
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("%+v", err)
	}
}
```
