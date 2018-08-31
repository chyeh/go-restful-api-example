package main

import (
	"fmt"
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

type ginHTTPServerConfig struct {
	host string
	port string
}

type ginHTTPServer struct {
	router  *gin.Engine
	address string
}

func newGinHTTPServer(cfg ginHTTPServerConfig) *ginHTTPServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	setRoutes(router)
	return &ginHTTPServer{
		router:  router,
		address: net.JoinHostPort(cfg.host, cfg.port),
	}
}

func (s *ginHTTPServer) run() {
	fmt.Println("starting http service on:", s.address)
	if err := s.router.Run(s.address); err != nil {
		fmt.Fprintln(os.Stderr, "error starting http service:", err)
		os.Exit(1)
	}
}

func setRoutes(router *gin.Engine) {
	router.GET("/recipes", listRecipe)
}
