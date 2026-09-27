package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 健康检查接口，k8s探针使用
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// hello业务接口
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "hello go k3s demo",
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 启动goroutine运行http服务
	go func() {
		log.Printf("server start listen on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen serve err: %v", err)
		}
	}()

	// 捕获退出信号 SIGINT(Ctrl‑C) SIGTERM(k8s pod删除发出的信号)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("received shutdown signal, start graceful shutdown ...")

	// 设置5s优雅关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server gracefully exited")
}