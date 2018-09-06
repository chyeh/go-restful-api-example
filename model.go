package main

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v9"
	null "gopkg.in/guregu/null.v3"
)

var validate = func() *validator.Validate {
	v := validator.New()
	v.RegisterCustomTypeFunc(
		func(field reflect.Value) interface{} {
			return field.Interface().(null.String).String
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
	Name         null.String `json:"name" db:"r_name" validate:"required"`
	PrepareTime  null.Int    `json:"prepare_time" db:"r_prep_time" validate:"omitempty,gt=0"`
	Difficulty   null.Int    `json:"difficulty" db:"r_difficulty" validate:"omitempty,min=1,max=3"`
	IsVegetarian null.Bool   `json:"is_vegetarian" db:"r_vegetarian" validate:"required"`
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
	Rating null.Int `json:"rating" validate:"required,min=1,max=5"`
}

func (a *PostRateRecipeArg) updateRatedRecipe(r *Recipe) {
	if r == nil {
		return
	}
	r.Rating.Float64 = ((r.Rating.Float64 * float64(r.RatedNum.Int64)) + float64(a.Rating.Int64)) / float64(r.RatedNum.Int64+1)
	r.RatedNum.Int64++
}

type filter struct {
	isSet          bool
	name           null.String
	prepTimeFrom   null.Int
	prepTimeTo     null.Int
	difficultyFrom null.Int
	difficultyTo   null.Int
	isVegetarian   null.Bool
}

func newFilter(c *gin.Context) *filter {
	var f filter
	if v, ok := c.GetQuery("name"); ok {
		f.name = null.StringFrom(v)
		f.isSet = true
	}
	if v, ok := c.GetQuery("prepare_time_from"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		f.prepTimeFrom = null.IntFrom(n)
		f.isSet = true
	}
	if v, ok := c.GetQuery("prepare_time_to"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		f.prepTimeTo = null.IntFrom(n)
		f.isSet = true
	}
	if v, ok := c.GetQuery("difficulty_from"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		f.difficultyFrom = null.IntFrom(n)
		f.isSet = true
	}
	if v, ok := c.GetQuery("difficulty_to"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		f.difficultyTo = null.IntFrom(n)
		f.isSet = true
	}
	if v, ok := c.GetQuery("is_vegetarian"); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		f.isVegetarian = null.BoolFrom(b)
		f.isSet = true
	}
	return &f
}

func (f *filter) whereClause() string {
	if !f.isSet {
		return ""
	}
	clause := " WHERE "
	if f.name.Valid {
		clause += (" r_name LIKE '%" + f.name.String + "%' AND ")
	}
	if f.prepTimeFrom.Valid {
		clause += (" r_prep_time >= " + strconv.Itoa(int(f.prepTimeFrom.Int64)) + " AND ")
	}
	if f.prepTimeTo.Valid {
		clause += (" r_prep_time <= " + strconv.Itoa(int(f.prepTimeTo.Int64)) + " AND ")
	}
	if f.difficultyFrom.Valid {
		clause += (" r_difficulty >= " + strconv.Itoa(int(f.difficultyFrom.Int64)) + " AND ")
	}
	if f.difficultyTo.Valid {
		clause += (" r_difficulty <= " + strconv.Itoa(int(f.difficultyTo.Int64)) + " AND ")
	}
	if f.isVegetarian.Valid {
		clause += (" r_vegetarian = " + strconv.FormatBool(f.isVegetarian.Bool) + " AND ")
	}
	clause += " true "
	return clause
}
