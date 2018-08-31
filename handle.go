package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	null "gopkg.in/guregu/null.v3"
)

type Recipe struct {
	ID           int      `json:"id" db:"r_id"`
	Name         string   `json:"name" db:"r_name"`
	PrepareTime  null.Int `json:"prepare_time" db:"r_prep_time"`
	Difficulty   null.Int `json:"difficulty" db:"r_difficulty"`
	IsVegetarian bool     `json:"is_vegetarian" db:"r_vegetarian"`
}

func listRecipe(c *gin.Context) {
	connStr := "postgres://hellofresh:hellofresh@localhost:5432/hellofresh?sslmode=disable"
	db := sqlx.MustConnect("postgres", connStr)
	var res []Recipe
	err := db.Select(&res, "SELECT * FROM recipe")
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, res)
}
