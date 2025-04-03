package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// RunHTTPServer starts an HTTP server.
func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		// TODO 加入告警日志
	}
	RunHTTPServerOnAddr(addr, wrapper)
}

// RunHTTPServerOnAddr starts an HTTP server on a specified address.
func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	wrapper(apiRouter)

	apiRouter.Group("/api")
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}
