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

type datastore interface {
	listRecipes() []*Recipe
	addRecipe(*PostRecipeArg) *Recipe
	getRecipeByID(int) *Recipe
}

type sqlxPostgreSQL struct {
	sqlxDB *sqlx.DB
}

func newSqlxPostgreSQL(connectionString string) *sqlxPostgreSQL {
	return &sqlxPostgreSQL{
		sqlxDB: sqlx.MustConnect("postgres", connectionString),
	}
}

func (d *sqlxPostgreSQL) close() {
	if err := d.sqlxDB.Close(); err != nil {
		panic(err)
	}
}

func (d *sqlxPostgreSQL) listRecipes() []*Recipe {
	res := make([]*Recipe, 0)
	if err := d.sqlxDB.Select(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian FROM recipe
	`); err != nil {
		panic(err)
	}
	return res
}

func (d *sqlxPostgreSQL) addRecipe(arg *PostRecipeArg) *Recipe {
	var res Recipe
	tx := d.sqlxDB.MustBegin()
	tx.MustExec(`
	INSERT INTO recipe(r_name, r_prep_time, r_difficulty, r_vegetarian)
	VALUES ($1, $2, $3, $4)
	`, arg.Name, arg.PrepareTime, arg.Difficulty, arg.IsVegetarian)
	tx.Get(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian FROM recipe
	WHERE r_id = (
		SELECT currval(pg_get_serial_sequence('recipe','r_id'))
	)
	`)
	tx.Commit()
	return &res
}

func (d *sqlxPostgreSQL) getRecipeByID(id int) *Recipe {
	var res Recipe
	if err := d.sqlxDB.Get(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian FROM recipe
	WHERE r_id = $1
	`, id); err != nil {
		return nil
	}
	return &res
}
