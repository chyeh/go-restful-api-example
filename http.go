package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type ginHTTPServer struct {
	*http.Server
	router *gin.Engine
}

func newGinHTTPServer() *ginHTTPServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(buildPanicProcessor(defaultPanicProcessor))
	return &ginHTTPServer{
		&http.Server{Handler: router},
		router,
	}
}

func (s *ginHTTPServer) run(address string) {
	s.Addr = address
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, "error starting http service:", err)
		os.Exit(1)
	}
}

type panicProcessor func(c *gin.Context, panic interface{})

func buildPanicProcessor(processor panicProcessor) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}

			processor(c, p)
		}()
		c.Next()
	}
}

func defaultPanicProcessor(c *gin.Context, panicObj interface{}) {
	httprequest, _ := httputil.DumpRequest(c.Request, false)
	msg := fmt.Sprintf("[Panic] %s\n%s\n%s\n", time.Now().Format(time.RFC3339), httprequest, panicObj)
	c.String(http.StatusInternalServerError, msg)
}
