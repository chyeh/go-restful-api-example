package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

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
	s.httpServer.router.GET("/recipe/:id", s.getRecipe)
	s.httpServer.router.PUT("/recipe/:id", s.putRecipe)
	s.httpServer.router.DELETE("/recipe/:id", s.deleteRecipe)
}

func (s *apiServer) getRecipes(c *gin.Context) {
	res := s.datastore.listRecipes(newFilter(c))
	c.JSON(200, res)
}

func (s *apiServer) postRecipes(c *gin.Context) {
	arg := &PostRecipeArg{}
	if err := c.BindJSON(arg); err == nil {
		arg.validate()
		res := s.datastore.addRecipe(arg)
		c.JSON(200, res)
	}
}

func (s *apiServer) getRecipe(c *gin.Context) {
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(404)
	} else if res := s.datastore.getRecipeByID(recipeID); res != nil {
		c.JSON(200, res)
	}
	c.AbortWithStatus(404)
}

func (s *apiServer) putRecipe(c *gin.Context) {
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	arg := &PutRecipeArg{}
	if err := c.BindJSON(arg); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token := c.GetHeader("Authorization")
	if recipe := s.datastore.updateAndGetRecipeByCredential(arg, recipeID, token); recipe != nil {
		c.JSON(200, recipe)
		return
	}

	c.AbortWithStatus(404)
}

func (s *apiServer) deleteRecipe(c *gin.Context) {
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	token := c.GetHeader("Authorization")
	if recipe := s.datastore.deleteRecipeByID(recipeID, token); recipe != nil {
		c.JSON(200, recipe)
		return
	}

	c.AbortWithStatus(404)
}
