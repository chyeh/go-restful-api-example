package main

import (
	null "gopkg.in/guregu/null.v3"
)

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
