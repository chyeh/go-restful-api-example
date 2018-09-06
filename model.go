package main

import (
	"fmt"
	"reflect"
	"strconv"

	validator "gopkg.in/go-playground/validator.v9"
	null "gopkg.in/guregu/null.v3"
)

var validate = func() *validator.Validate {
	v := validator.New()
	v.RegisterCustomTypeFunc(
		func(field reflect.Value) interface{} {
			return field.Interface().(null.String).Ptr()
		},
		null.String{},
	)
	v.RegisterCustomTypeFunc(
		func(field reflect.Value) interface{} {
			return field.Interface().(null.Int).Ptr()
		},
		null.Int{},
	)
	v.RegisterCustomTypeFunc(
		func(field reflect.Value) interface{} {
			return field.Interface().(null.Bool).Ptr()
		},
		null.Bool{},
	)
	v.RegisterCustomTypeFunc(
		func(field reflect.Value) interface{} {
			return field.Interface().(null.Float).Ptr()
		},
		null.Float{},
	)
	return v
}()

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
	Name         null.String `json:"name" db:"r_name" validate:"required,gt=0"`
	PrepareTime  null.Int    `json:"prepare_time" db:"r_prep_time" validate:"omitempty,gt=0"`
	Difficulty   null.Int    `json:"difficulty" db:"r_difficulty" validate:"omitempty,min=1,max=3"`
	IsVegetarian null.Bool   `json:"is_vegetarian" db:"r_vegetarian" validate:"required"`
}

type PutRecipeArg struct {
	Name         null.String `json:"name" validate:"omitempty,gt=0"`
	PrepareTime  null.Int    `json:"prepare_time" validate:"omitempty,gt=0"`
	Difficulty   null.Int    `json:"difficulty" validate:"omitempty,min=1,max=3"`
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
	Rating null.Int `json:"rating" validate:"required,min=1,max=5"`
}

func (a *PostRateRecipeArg) updateRatedRecipe(r *Recipe) {
	if r == nil {
		return
	}
	r.Rating.Float64 = ((r.Rating.Float64 * float64(r.RatedNum.Int64)) + float64(a.Rating.Int64)) / float64(r.RatedNum.Int64+1)
	r.RatedNum.Int64++
}

type ListFilter struct {
	Name           string `form:"name"`
	PrepTimeFrom   int    `form:"prepare_time_from"`
	PrepTimeTo     int    `form:"prepare_time_to"`
	DifficultyFrom int    `form:"difficulty_from"`
	DifficultyTo   int    `form:"difficulty_to"`
	IsVegetarian   string `form:"is_vegetarian"`
}

func (f *ListFilter) whereClause() string {
	clause := " WHERE "
	clause += (" r_name LIKE '%" + f.Name + "%' AND ")
	if f.PrepTimeFrom != 0 {
		clause += fmt.Sprintf(" r_prep_time >= %d AND ", f.PrepTimeFrom)
	}
	if f.PrepTimeTo != 0 {
		clause += fmt.Sprintf(" r_prep_time <= %d AND ", f.PrepTimeTo)
	}
	if f.DifficultyFrom != 0 {
		clause += fmt.Sprintf(" r_difficulty >= %d AND ", f.DifficultyFrom)
	}
	if f.DifficultyTo != 0 {
		clause += fmt.Sprintf(" r_difficulty <= %d AND ", f.DifficultyTo)
	}
	if f.IsVegetarian != "" {
		b, err := strconv.ParseBool(f.IsVegetarian)
		if err != nil {
			panic(err)
		}
		clause += fmt.Sprintf(" r_vegetarian = %v ", b)
	} else {
		clause += " true "
	}
	return clause
}
