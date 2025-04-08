package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	setMiddleware(apiRouter)
	wrapper(apiRouter)

	apiRouter.Group("/api")
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}

func setMiddleware(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(otelgin.Middleware("default_server"))
}
