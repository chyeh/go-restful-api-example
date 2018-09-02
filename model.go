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
