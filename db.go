package main

import (
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

type sqlxPostgreSQL struct {
	sqlxDB *sqlx.DB
}

func newSqlxPostgreSQL(connectionString string) *sqlxPostgreSQL {
	return &sqlxPostgreSQL{
		sqlxDB: sqlx.MustConnect("postgres", connectionString),
	}
}

func (d *sqlxPostgreSQL) listRecipes() []Recipe {
	var res []Recipe
	if err := d.sqlxDB.Select(&res, "SELECT * FROM recipe"); err != nil {
		panic(err)
	}
	return res
}
