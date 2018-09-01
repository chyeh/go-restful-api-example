package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type ginHTTPServer struct {
	router *gin.Engine
}

func newGinHTTPServer() *ginHTTPServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	return &ginHTTPServer{
		router: router,
	}
}

func (s *ginHTTPServer) run(address string) {
	if err := s.router.Run(address); err != nil {
		fmt.Fprintln(os.Stderr, "error starting http service:", err)
		os.Exit(1)
	}
}
