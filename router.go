/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/rosbit/mgin"
	"chatGPT-gateway/conf"
	"chatGPT-gateway/rest"
	"net/http"
	"fmt"
	"os"
	"log"
	"syscall"
	"os/signal"
)

// 设置路由，进入服务状态
func StartService() error {
	initService()
	serviceConf := &conf.ServiceConf

	api := mgin.NewMgin(mgin.WithLogger("chatGPT-gateway"))

	ce := &serviceConf.CommonEndpoints
	// health check
	api.GET(ce.HealthCheck, func(c *mgin.Context) {
		c.String(http.StatusOK, "OK\n")
	})
	api.POST(ce.Chat, rest.Chat)
	api.POST(ce.Image, rest.CreateImage)
	api.POST(ce.ChatCompletions, rest.ChatCompletions)

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	log.Printf("I am listening at %s...\n", listenParam)
	return http.ListenAndServe(listenParam, api)
}

func initService() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGQUIT)
	go func() {
		for range c {
			log.Println("I will exit in a while")
			os.Exit(0)
		}
	}()
}

