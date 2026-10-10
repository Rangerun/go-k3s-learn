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
	appVer := os.Getenv("APP_VER")
	appMsg := os.Getenv("APP_MSG")
	dbPass := os.Getenv("DB_PASSWORD")
	// 健康检查接口，k8s探针使用
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"version": "v2",
			"db_pass":  dbPass, //练习演示用！线上禁止返回密钥
		})
	})

	//原来v1返回: hello go k3s demo
	//v2修改接口返回
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ver": appVer,
			"msg": appMsg,
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