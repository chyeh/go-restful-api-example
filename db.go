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
	listRecipes(*filter) []*Recipe
	addRecipe(*PostRecipeArg) *Recipe
	getRecipeByID(int) *Recipe
	updateRecipe(*Recipe)
	deleteRecipeByID(id int)
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

func (d *sqlxPostgreSQL) listRecipes(f *filter) []*Recipe {
	res := make([]*Recipe, 0)
	if err := d.sqlxDB.Select(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian FROM recipe
	`+f.whereClause()); err != nil {
		panic(err)
	}
	return res
}

func (d *sqlxPostgreSQL) addRecipe(arg *PostRecipeArg) *Recipe {
	var res Recipe
	tx := d.sqlxDB.MustBegin()
	if _, err := tx.NamedExec(`
	INSERT INTO recipe(r_name, r_prep_time, r_difficulty, r_vegetarian)
	VALUES (:r_name, :r_prep_time, :r_difficulty, :r_vegetarian)
	`, arg); err != nil {
		panic(err)
	}
	if err := tx.Get(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian FROM recipe
	WHERE r_id = (
		SELECT currval(pg_get_serial_sequence('recipe','r_id'))
	)
	`); err != nil {
		panic(err)
	}
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

func (d *sqlxPostgreSQL) updateRecipe(arg *Recipe) {
	if _, err := d.sqlxDB.NamedExec(`
	UPDATE recipe
	SET	r_name = :r_name,
		r_prep_time = :r_prep_time,
		r_difficulty = :r_difficulty,
		r_vegetarian = :r_vegetarian
	WHERE r_id = :r_id
	`, arg); err != nil {
		panic(err)
	}
}

func (d *sqlxPostgreSQL) deleteRecipeByID(id int) {
	d.sqlxDB.MustExec(`
	DELETE FROM recipe
	WHERE r_id = $1
	`, id)
}
