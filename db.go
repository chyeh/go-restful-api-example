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
	deleteRecipeByID(id int)
	updateAndGetRecipeByCredential(*PutRecipeArg, int, string) *Recipe
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
	if f == nil {
		panic("nil filter not allowed")
	}
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

func (d *sqlxPostgreSQL) updateAndGetRecipeByCredential(arg *PutRecipeArg, id int, token string) *Recipe {
	var res Recipe
	tx := d.sqlxDB.MustBegin()
	if err := d.sqlxDB.Get(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian FROM recipe
	WHERE r_id = $1
	`, id); err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		return nil
	}
	arg.overwriteRecipe(&res)
	slqResult := tx.MustExec(`
	UPDATE recipe
	SET	r_name = $1,
		r_prep_time = $2,
		r_difficulty = $3,
		r_vegetarian = $4
	WHERE r_id = (
		SELECT recipe.r_id
		FROM recipe
		INNER JOIN hellofresh_user_recipe
		ON recipe.r_id = hellofresh_user_recipe.hur_r_id
		WHERE recipe.r_id = $5 AND hellofresh_user_recipe.hur_hu_id= (
			SELECT hu_id FROM hellofresh_user
			WHERE hu_access_token = $6
		)
	)
	`, res.Name, res.PrepareTime, res.Difficulty, res.IsVegetarian, id, token)
	if cnt, err := slqResult.RowsAffected(); err != nil {
		panic(err)
	} else if cnt == 0 {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		return nil
	}
	tx.Commit()
	return &res
}

func (d *sqlxPostgreSQL) deleteRecipeByID(id int) {
	d.sqlxDB.MustExec(`
	DELETE FROM recipe
	WHERE r_id = $1
	`, id)
}
