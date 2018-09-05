package main

import (
	null "gopkg.in/guregu/null.v3"
)

type Recipe struct {
	ID           int        `json:"id" db:"r_id"`
	Name         string     `json:"name" db:"r_name"`
	PrepareTime  null.Int   `json:"prepare_time" db:"r_prep_time"`
	Difficulty   null.Int   `json:"difficulty" db:"r_difficulty"`
	IsVegetarian bool       `json:"is_vegetarian" db:"r_vegetarian"`
	Rating       null.Float `json:"rating" db:"r_rating"`
	RatedNum     null.Int   `json:"rated_num" db:"r_rated_num"`
}

type PostRecipeArg struct {
	Name         null.String `json:"name" db:"r_name"`
	PrepareTime  null.Int    `json:"prepare_time" db:"r_prep_time"`
	Difficulty   null.Int    `json:"difficulty" db:"r_difficulty"`
	IsVegetarian null.Bool   `json:"is_vegetarian" db:"r_vegetarian"`
}

func (a *PostRecipeArg) validate() {
	if !a.Name.Valid {
		panic("field 'Name' not valid")
	}
	if !a.IsVegetarian.Valid {
		panic("field 'IsVegetarian' not valid")
	}
}

type PutRecipeArg struct {
	Name         null.String `json:"name"`
	PrepareTime  null.Int    `json:"prepare_time"`
	Difficulty   null.Int    `json:"difficulty"`
	IsVegetarian null.Bool   `json:"is_vegetarian"`
}

func (a *PutRecipeArg) overwriteRecipe(r *Recipe) {
	if r == nil {
		return
	}
	if a.Name.Valid {
		r.Name = a.Name.String
	}
	if a.PrepareTime.Valid {
		r.PrepareTime = a.PrepareTime
	}
	if a.Difficulty.Valid {
		r.Difficulty = a.Difficulty
	}
	if a.IsVegetarian.Valid {
		r.IsVegetarian = a.IsVegetarian.Bool
	}
}

type PostRateRecipeArg struct {
	Rating null.Int `json:"rating"`
}

func (a *PostRateRecipeArg) validate() {
	if !a.Rating.Valid {
		panic("field 'Rating' not valid")
	}
}

func (a *PostRateRecipeArg) updateRecipe(r *Recipe) {
	if r == nil {
		return
	}
	r.Rating.Float64 = ((r.Rating.Float64 * float64(r.RatedNum.Int64)) + float64(a.Rating.Int64)) / float64(r.RatedNum.Int64+1)
}
