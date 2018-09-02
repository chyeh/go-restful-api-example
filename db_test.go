package main

import (
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	null "gopkg.in/guregu/null.v3"
)

var _ = Describe("Testing database object", func() {
	BeforeEach(func() {
		testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@10.20.30.50:5432/?sslmode=disable")
		defer testDB.close()

		testDB.sqlxDB.MustExec(`
		DROP DATABASE IF EXISTS test_hellofresh
		`)
		testDB.sqlxDB.MustExec(`
		CREATE DATABASE test_hellofresh
		`)
	})

	AfterEach(func() {
		testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@10.20.30.50:5432/?sslmode=disable")
		defer testDB.close()

		testDB.sqlxDB.MustExec(`
		DROP DATABASE test_hellofresh
		`)
	})
	Context("listing recipes", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@10.20.30.50:5432/test_hellofresh?sslmode=disable")
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
				r_vegetarian BOOLEAN NOT NULL DEFAULT false
			)
			`)
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@10.20.30.50:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
			DROP TABLE recipe
			`)
		})
		It("lists empty table", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@10.20.30.50:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			actual := testDB.listRecipes()
			Expect(actual).NotTo(BeNil())
			Expect(actual).To(HaveLen(0))
		})

		It("lists non-empty table", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@10.20.30.50:5432/test_hellofresh?sslmode=disable")
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
})
