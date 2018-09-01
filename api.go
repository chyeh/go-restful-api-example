package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	null "gopkg.in/guregu/null.v3"
)

type apiServerConfig struct {
	host             string
	port             string
	connectionString string
}

type apiServer struct {
	httpServer     *ginHTTPServer
	sqlxPostgreSQL *sqlx.DB
}

func newAPIServer(cfg apiServerConfig) *apiServer {
	httpServer := newGinHTTPServer(ginHTTPServerConfig{
		host: cfg.host,
		port: cfg.port,
	})
	apiServer := &apiServer{
		httpServer:     httpServer,
		sqlxPostgreSQL: sqlx.MustConnect("postgres", cfg.connectionString),
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

type Recipe struct {
	ID           int      `json:"id" db:"r_id"`
	Name         string   `json:"name" db:"r_name"`
	PrepareTime  null.Int `json:"prepare_time" db:"r_prep_time"`
	Difficulty   null.Int `json:"difficulty" db:"r_difficulty"`
	IsVegetarian bool     `json:"is_vegetarian" db:"r_vegetarian"`
}

func (s *apiServer) GetRecipes(c *gin.Context) {
	var res []Recipe
	err := s.sqlxPostgreSQL.Select(&res, "SELECT * FROM recipe")
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, res)
}
