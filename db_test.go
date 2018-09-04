package main

import (
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	null "gopkg.in/guregu/null.v3"
)

const (
	testRecipeTableSchema = `
	CREATE TABLE recipe(
		r_id SERIAL PRIMARY KEY,
		r_name VARCHAR(512) NOT NULL,
		r_prep_time SMALLINT,
		r_difficulty SMALLINT,
		r_vegetarian BOOLEAN NOT NULL
	)
	`
	testHellofreshUserTableSchema = `
	CREATE TABLE hellofresh_user(
		hu_id SERIAL PRIMARY KEY,
		hu_account VARCHAR(32) NOT NULL UNIQUE,
		hu_access_token VARCHAR(32) NOT NULL UNIQUE
	)
	`
	testHellofreshUserRecipeTableSchema = `
	CREATE TABLE hellofresh_user_recipe(
		hur_hu_id INTEGER,
		hur_r_id INTEGER,
		CONSTRAINT pk_hellofresh_user_recipe PRIMARY KEY(hur_hu_id, hur_r_id),
		CONSTRAINT fk_hellofresh_user_recipe__hellofresh_user FOREIGN KEY
			(hur_hu_id) REFERENCES hellofresh_user(hu_id)
			ON DELETE CASCADE
			ON UPDATE RESTRICT,
		CONSTRAINT fk_hellofresh_user_recipe__recipe FOREIGN KEY
			(hur_r_id) REFERENCES recipe(r_id)
			ON DELETE CASCADE
			ON UPDATE RESTRICT
	)
	`
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
			testDB.sqlxDB.MustExec(testRecipeTableSchema)
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

			actual := testDB.listRecipes(&filter{})
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
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(3))
		})
		It("lists non-empty table with filters", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name1", true),
				PrepareTime:  null.NewInt(15, true),
				Difficulty:   null.NewInt(2, true),
				IsVegetarian: null.NewBool(false, true),
			})
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name2", true),
				PrepareTime:  null.NewInt(20, true),
				Difficulty:   null.NewInt(1, true),
				IsVegetarian: null.NewBool(false, true),
			})
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name3", true),
				PrepareTime:  null.NewInt(50, true),
				Difficulty:   null.NewInt(3, true),
				IsVegetarian: null.NewBool(true, true),
			})
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name4", true),
				PrepareTime:  null.NewInt(60, true),
				Difficulty:   null.NewInt(5, true),
				IsVegetarian: null.NewBool(false, true),
			})
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name5", true),
				PrepareTime:  null.NewInt(70, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: null.NewBool(true, true),
			})
			Expect(testDB.listRecipes(&filter{
				name:  null.NewString("name", true),
				isSet: true,
			})).To(HaveLen(5))
			Expect(testDB.listRecipes(&filter{
				name:  null.NewString("5", true),
				isSet: true,
			})).To(HaveLen(1))
			Expect(testDB.listRecipes(&filter{
				name:  null.NewString("x", true),
				isSet: true,
			})).To(HaveLen(0))
			Expect(testDB.listRecipes(&filter{
				isVegetarian: null.NewBool(true, true),
				isSet:        true,
			})).To(HaveLen(2))
			Expect(testDB.listRecipes(&filter{
				difficultyTo: null.NewInt(3, true),
				isVegetarian: null.NewBool(true, true),
				isSet:        true,
			})).To(HaveLen(1))
			Expect(testDB.listRecipes(&filter{
				prepTimeTo: null.NewInt(60, true),
				isSet:      true,
			})).To(HaveLen(4))
			Expect(testDB.listRecipes(&filter{
				prepTimeFrom: null.NewInt(20, true),
				prepTimeTo:   null.NewInt(60, true),
				isSet:        true,
			})).To(HaveLen(3))
			Expect(testDB.listRecipes(&filter{
				prepTimeFrom: null.NewInt(20, true),
				prepTimeTo:   null.NewInt(60, true),
				difficultyTo: null.NewInt(3, true),
				isSet:        true,
			})).To(HaveLen(2))
			Expect(testDB.listRecipes(&filter{
				prepTimeFrom:   null.NewInt(20, true),
				prepTimeTo:     null.NewInt(60, true),
				difficultyFrom: null.NewInt(2, true),
				difficultyTo:   null.NewInt(4, true),
				isSet:          true,
			})).To(HaveLen(1))
			Expect(testDB.listRecipes(&filter{
				prepTimeFrom:   null.NewInt(20, true),
				prepTimeTo:     null.NewInt(60, true),
				difficultyFrom: null.NewInt(2, true),
				difficultyTo:   null.NewInt(4, true),
				isVegetarian:   null.NewBool(false, true),
				isSet:          true,
			})).To(HaveLen(0))
		})
	})
	Context("adding a new recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
				DROP TABLE IF EXISTS recipe
				`)
			testDB.sqlxDB.MustExec(testRecipeTableSchema)
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

			Expect(testDB.listRecipes(&filter{})).To(HaveLen(0))
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name1", true),
				IsVegetarian: null.NewBool(false, true),
			})
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(1))

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
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(2))
		})
	})
	Context("updating a recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
				DROP TABLE IF EXISTS recipe
				`)
			testDB.sqlxDB.MustExec(testRecipeTableSchema)
			testDB.sqlxDB.MustExec(`
				DROP TABLE IF EXISTS hellofresh_user
				`)
			testDB.sqlxDB.MustExec(testHellofreshUserTableSchema)
			testDB.sqlxDB.MustExec(`
			DROP TABLE IF EXISTS hellofresh_user_recipe
			`)
			testDB.sqlxDB.MustExec(testHellofreshUserRecipeTableSchema)

			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name1", true),
				PrepareTime:  null.NewInt(2, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: null.NewBool(false, true),
			})
			testDB.sqlxDB.MustExec(`
			INSERT INTO hellofresh_user(hu_account, hu_access_token)
			VALUES
			('foo', 'faketoken')
			`)
			testDB.sqlxDB.MustExec(`
			INSERT INTO hellofresh_user_recipe(hur_hu_id, hur_r_id)
			VALUES
			(1,1)
			`)
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
			DROP TABLE hellofresh_user_recipe
			`)
			testDB.sqlxDB.MustExec(`
			DROP TABLE hellofresh_user
			`)
			testDB.sqlxDB.MustExec(`
			DROP TABLE recipe
			`)
		})
		It("updates a existent record in the recipe table", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			actual := testDB.updateAndGetRecipeByCredential(&PutRecipeArg{
				Name:         null.NewString("name1_updated", true),
				PrepareTime:  null.NewInt(3, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: null.NewBool(false, true),
			}, 1, "faketoken")

			Expect(actual.Name).To(Equal("name1_updated"))
			Expect(actual.PrepareTime.Int64).To(Equal(int64(3)))
			Expect(actual.Difficulty.Int64).To(Equal(int64(4)))
			Expect(actual.IsVegetarian).To(BeFalse())
		})
		It("does nothing if the recipe doesn't exist", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			actual := testDB.updateAndGetRecipeByCredential(&PutRecipeArg{
				Name:         null.NewString("name1_updated", true),
				PrepareTime:  null.NewInt(3, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: null.NewBool(false, true),
			}, 2, "faketoken")

			Expect(actual).To(BeNil())
		})
		It("does nothing if the access to the recipe is not authorized", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			actual := testDB.updateAndGetRecipeByCredential(&PutRecipeArg{
				Name:         null.NewString("name1_updated", true),
				PrepareTime:  null.NewInt(3, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: null.NewBool(false, true),
			}, 1, "failed_faketoken")

			Expect(actual).To(BeNil())
		})
	})
	Context("deleting a recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
				DROP TABLE IF EXISTS recipe
				`)
			testDB.sqlxDB.MustExec(testRecipeTableSchema)
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.sqlxDB.MustExec(`
				DROP TABLE recipe
				`)
		})
		It("deletes a existent record in the recipe table", func() {
			testDB := newSqlxPostgreSQL("postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable")
			defer testDB.close()

			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name1", true),
				PrepareTime:  null.NewInt(1, true),
				Difficulty:   null.NewInt(2, true),
				IsVegetarian: null.NewBool(false, true),
			})
			testDB.addRecipe(&PostRecipeArg{
				Name:         null.NewString("name2", true),
				PrepareTime:  null.NewInt(2, true),
				Difficulty:   null.NewInt(4, true),
				IsVegetarian: null.NewBool(true, true),
			})
			testDB.deleteRecipeByID(1)
			testDB.getRecipeByID(1)
			Expect(testDB.getRecipeByID(1)).To(BeNil())
			Expect(testDB.getRecipeByID(2)).NotTo(BeNil())
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(1))
			testDB.deleteRecipeByID(2)
			Expect(testDB.getRecipeByID(2)).To(BeNil())
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(0))
		})
	})
})
