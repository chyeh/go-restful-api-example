package main

import (
	"github.com/gin-gonic/gin"
)

type apiServerConfig struct {
	host             string
	port             string
	connectionString string
}

type apiServer struct {
	httpServer *ginHTTPServer
	datastore  *sqlxPostgreSQL
}

func newAPIServer(cfg apiServerConfig) *apiServer {
	httpServer := newGinHTTPServer(ginHTTPServerConfig{
		host: cfg.host,
		port: cfg.port,
	})
	apiServer := &apiServer{
		httpServer: httpServer,
		datastore:  newSqlxPostgreSQL(cfg.connectionString),
	}
	apiServer.routes()
	return apiServer
}

func (s *apiServer) run() {
	s.httpServer.run()
}

func (s *apiServer) routes() {
	s.httpServer.router.GET("/recipes", s.GetRecipes)
}

func (s *apiServer) GetRecipes(c *gin.Context) {
	res := s.datastore.listRecipes()
	c.JSON(200, res)
}
