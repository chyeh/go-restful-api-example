package main

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

type apiServerConfig struct {
	host             string
	port             string
	connectionString string
}

type apiServer struct {
	httpServer *ginHTTPServer
	address    string
	datastore  datastore
}

func newAPIServer(cfg apiServerConfig) *apiServer {
	httpServer := newGinHTTPServer()
	apiServer := &apiServer{
		httpServer: httpServer,
		address:    net.JoinHostPort(cfg.host, cfg.port),
		datastore:  newSqlxPostgreSQL(cfg.connectionString),
	}
	apiServer.routes()
	return apiServer
}

func (s *apiServer) run() {
	fmt.Println("starting http service on:", s.address)
	s.httpServer.run(s.address)
}

func (s *apiServer) routes() {
	s.httpServer.router.GET("/recipes", s.getRecipes)
	s.httpServer.router.POST("/recipes", s.postRecipes)
}

func (s *apiServer) getRecipes(c *gin.Context) {
	res := s.datastore.listRecipes()
	c.JSON(200, res)
}

func (s *apiServer) postRecipes(c *gin.Context) {
	arg := PostRecipeArg{}
	if err := c.BindJSON(&arg); err == nil {
		res := s.datastore.addRecipe(arg)
		c.JSON(200, res)
	}
}
