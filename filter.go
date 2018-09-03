package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	null "gopkg.in/guregu/null.v3"
)

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
		f.name = null.NewString(v, true)
		f.isSet = true
	}
	if v, ok := c.GetQuery("prepare_time_from"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		f.prepTimeFrom = null.NewInt(n, true)
		f.isSet = true
	}
	if v, ok := c.GetQuery("prepare_time_to"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		f.prepTimeTo = null.NewInt(n, true)
		f.isSet = true
	}
	if v, ok := c.GetQuery("difficulty_from"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		f.difficultyFrom = null.NewInt(n, true)
		f.isSet = true
	}
	if v, ok := c.GetQuery("difficulty_to"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		f.difficultyTo = null.NewInt(n, true)
		f.isSet = true
	}
	if v, ok := c.GetQuery("is_vegetarian"); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		f.isVegetarian = null.NewBool(b, true)
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
