package main

import (
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	null "gopkg.in/guregu/null.v3"
)

var _ = Describe("Testing database object", func() {
	BeforeEach(func() {
		testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/?sslmode=disable")
		defer testDB.close()

		testDB.sqlxDB.MustExec(`
		DROP DATABASE IF EXISTS test_hellofresh
		`)
		testDB.sqlxDB.MustExec(`
		CREATE DATABASE test_hellofresh
		`)
	})
	AfterEach(func() {
		testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/?sslmode=disable")
		defer testDB.close()

		testDB.sqlxDB.MustExec(`
		DROP DATABASE test_hellofresh
		`)
	})
	Context("listing recipes", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
			DROP TABLE IF EXISTS recipe
			`)
			testDB.sqlxDB.MustExec(`
			CREATE TABLE recipe(
				r_id SERIAL PRIMARY KEY,
				r_name VARCHAR(512) NOT NULL,
				r_prep_time SMALLINT,
				r_difficulty SMALLINT,
				r_vegetarian BOOLEAN NOT NULL
			)
			`)
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
			DROP TABLE recipe
			`)
		})
		It("lists empty table", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			actual := testDB.listRecipes()
			Expect(actual).NotTo(BeNil())
			Expect(actual).To(HaveLen(0))
		})
		It("lists non-empty table", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name1", true),
				IsVegetarian: null.NewBool(false, true),
			})
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name2", true),
				IsVegetarian: null.NewBool(false, true),
			})
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name3", true),
				IsVegetarian: null.NewBool(false, true),
			})
			Expect(testDB.listRecipes()).To(HaveLen(3))
		})
	})
	Context("adding a new recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
			DROP TABLE IF EXISTS recipe
			`)
			testDB.sqlxDB.MustExec(`
			CREATE TABLE recipe(
				r_id SERIAL PRIMARY KEY,
				r_name VARCHAR(512) NOT NULL,
				r_prep_time SMALLINT,
				r_difficulty SMALLINT,
				r_vegetarian BOOLEAN NOT NULL
			)
			`)
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
			DROP TABLE recipe
			`)
		})
		It("adds a record in the recipe table and return the corresponding record", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			Expect(testDB.listRecipes()).To(HaveLen(0))
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name1", true),
				IsVegetarian: null.NewBool(false, true),
			})
			Expect(testDB.listRecipes()).To(HaveLen(1))

			actual := testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name2", true),
				PrepareTime:  null.NewInt(2, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: null.NewBool(false, true),
			})
			Expect(actual.ID).To(Equal(2))
			Expect(actual.Name).To(Equal("name2"))
			Expect(actual.PrepareTime.Int64).To(Equal(int64(2)))
			Expect(actual.Difficulty.Int64).To(Equal(int64(4)))
			Expect(actual.IsVegetarian).To(BeFalse())
			Expect(testDB.listRecipes()).To(HaveLen(2))
		})
	})
	Context("updating a recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
			DROP TABLE IF EXISTS recipe
			`)
			testDB.sqlxDB.MustExec(`
			CREATE TABLE recipe(
				r_id SERIAL PRIMARY KEY,
				r_name VARCHAR(512) NOT NULL,
				r_prep_time SMALLINT,
				r_difficulty SMALLINT,
				r_vegetarian BOOLEAN NOT NULL
			)
			`)
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
			DROP TABLE recipe
			`)
		})
		It("updates a existent record in the recipe table", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name2", true),
				PrepareTime:  null.NewInt(2, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: null.NewBool(false, true),
			})
			testDB.updateRecipe(&Recipe{
				ID:           1,
				Name:         "name2_updated",
				PrepareTime:  null.NewInt(3, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: false,
			})

			actual := testDB.getRecipeByID(1)
			Expect(actual.Name).To(Equal("name2_updated"))
			Expect(actual.PrepareTime.Int64).To(Equal(int64(3)))
			Expect(actual.Difficulty.Int64).To(Equal(int64(4)))
			Expect(actual.IsVegetarian).To(BeFalse())
		})

	})
})
