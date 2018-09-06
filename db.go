package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type datastore interface {
	listRecipes(*ListFilter, *paging) []*Recipe
	addRecipeByCredential(*PostRecipeArg, string) *Recipe
	getRecipeByID(int) *Recipe
	updateAndGetRecipeByCredential(*PutRecipeArg, int, string) *Recipe
	deleteAndGetRecipeByCredential(int, string) *Recipe
	rateAndGetRecipe(*PostRateRecipeArg, int) *Recipe
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

func (d *sqlxPostgreSQL) listRecipes(f *ListFilter, p *paging) []*Recipe {
	if f == nil {
		panic("nil filter not allowed")
	}
	res := make([]*Recipe, 0)
	if err := d.sqlxDB.Select(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian, r_rating, r_rated_num FROM recipe
	`+f.whereClause()+`ORDER BY r_id`+p.limitClause()+p.offsetClause()); err != nil {
		panic(err)
	}
	return res
}

func (d *sqlxPostgreSQL) addRecipeByCredential(arg *PostRecipeArg, token string) *Recipe {
	var res Recipe
	var userID int
	tx := d.sqlxDB.MustBegin()
	if err := d.sqlxDB.Get(&userID, `
	SELECT hu_id FROM hellofresh_user
	WHERE hu_access_token = $1
	`, token); err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		return nil
	}
	if _, err := tx.NamedExec(`
	INSERT INTO recipe(r_name, r_prep_time, r_difficulty, r_vegetarian)
	VALUES (:r_name, :r_prep_time, :r_difficulty, :r_vegetarian)
	`, arg); err != nil {
		panic(err)
	}
	if err := tx.Get(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian, r_rating, r_rated_num FROM recipe
	WHERE r_id = (
		SELECT currval(pg_get_serial_sequence('recipe','r_id'))
	)
	`); err != nil {
		panic(err)
	}
	tx.MustExec(`
	INSERT INTO hellofresh_user_recipe(hur_hu_id, hur_r_id)
	VALUES ($1, $2)
	`, userID, res.ID)
	tx.Commit()
	return &res
}

func (d *sqlxPostgreSQL) getRecipeByID(id int) *Recipe {
	var res Recipe
	if err := d.sqlxDB.Get(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian, r_rating, r_rated_num FROM recipe
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
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian, r_rating, r_rated_num FROM recipe
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

func (d *sqlxPostgreSQL) deleteAndGetRecipeByCredential(id int, token string) *Recipe {
	var res Recipe
	tx := d.sqlxDB.MustBegin()
	if err := d.sqlxDB.Get(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian, r_rating, r_rated_num FROM recipe
	WHERE r_id = $1
	`, id); err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		return nil
	}
	slqResult := d.sqlxDB.MustExec(`
	DELETE FROM recipe
	WHERE r_id = (
		SELECT recipe.r_id
		FROM recipe
		INNER JOIN hellofresh_user_recipe
		ON recipe.r_id = hellofresh_user_recipe.hur_r_id
		WHERE recipe.r_id = $1 AND hellofresh_user_recipe.hur_hu_id= (
			SELECT hu_id FROM hellofresh_user
			WHERE hu_access_token = $2
		)
	)
	`, id, token)
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

func (d *sqlxPostgreSQL) rateAndGetRecipe(arg *PostRateRecipeArg, id int) *Recipe {
	var res Recipe
	tx := d.sqlxDB.MustBegin()
	if err := d.sqlxDB.Get(&res, `
	SELECT r_id, r_name, r_prep_time, r_difficulty, r_vegetarian, r_rating, r_rated_num FROM recipe
	WHERE r_id = $1
	`, id); err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		return nil
	}
	arg.updateRatedRecipe(&res)
	slqResult := tx.MustExec(`
	UPDATE recipe
	SET	r_rating = ((r_rating*r_rated_num) + $1)/(r_rated_num + 1),
		r_rated_num = r_rated_num + 1
	WHERE r_id = $2
	`, arg.Rating, id)
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
