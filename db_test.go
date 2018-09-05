package main

import (
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	null "gopkg.in/guregu/null.v3"
)

const (
	testDBConnectionString             = "postgres://hellofresh:hellofresh@localhost:5432/?sslmode=disable"
	testDBConnectionStringWithDatabase = "postgres://hellofresh:hellofresh@localhost:5432/test_hellofresh?sslmode=disable"
	testRecipeTableSchema              = `
	CREATE TABLE recipe(
		r_id SERIAL PRIMARY KEY,
		r_name VARCHAR(512) NOT NULL,
		r_prep_time SMALLINT,
		r_difficulty SMALLINT,
		r_vegetarian BOOLEAN NOT NULL,
		r_rating REAL NOT NULL DEFAULT 0.0,
		r_rated_num INTEGER NOT NULL DEFAULT 0
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
		testDB := newSqlxPostgreSQL(testDBConnectionString)
		defer testDB.close()

		testDB.sqlxDB.MustExec(`
		DROP DATABASE IF EXISTS test_hellofresh
		`)
		testDB.sqlxDB.MustExec(`
		CREATE DATABASE test_hellofresh
		`)
	})
	AfterEach(func() {
		testDB := newSqlxPostgreSQL(testDBConnectionString)
		defer testDB.close()

		testDB.sqlxDB.MustExec(`
		DROP DATABASE test_hellofresh
		`)
	})
	Context("listing recipes", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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
			testDB.sqlxDB.MustExec(`
			INSERT INTO hellofresh_user(hu_account, hu_access_token)
			VALUES
			('foo', 'faketoken')
			`)
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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
		It("lists empty table", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			actual := testDB.listRecipes(&filter{})
			Expect(actual).NotTo(BeNil())
			Expect(actual).To(HaveLen(0))
		})
		It("lists non-empty table", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name1"),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name2"),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name3"),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(3))
		})
		It("lists non-empty table with filters", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name1"),
				PrepareTime:  null.IntFrom(15),
				Difficulty:   null.IntFrom(2),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name2"),
				PrepareTime:  null.IntFrom(20),
				Difficulty:   null.IntFrom(1),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name3"),
				PrepareTime:  null.IntFrom(50),
				Difficulty:   null.IntFrom(3),
				IsVegetarian: null.BoolFrom(true),
			}, "faketoken")
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name4"),
				PrepareTime:  null.IntFrom(60),
				Difficulty:   null.IntFrom(5),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name5"),
				PrepareTime:  null.IntFrom(70),
				Difficulty:   null.IntFrom(4),
				IsVegetarian: null.BoolFrom(true),
			}, "faketoken")
			Expect(testDB.listRecipes(&filter{
				name:  null.StringFrom("name"),
				isSet: true,
			})).To(HaveLen(5))
			Expect(testDB.listRecipes(&filter{
				name:  null.StringFrom("5"),
				isSet: true,
			})).To(HaveLen(1))
			Expect(testDB.listRecipes(&filter{
				name:  null.StringFrom("x"),
				isSet: true,
			})).To(HaveLen(0))
			Expect(testDB.listRecipes(&filter{
				isVegetarian: null.BoolFrom(true),
				isSet:        true,
			})).To(HaveLen(2))
			Expect(testDB.listRecipes(&filter{
				difficultyTo: null.IntFrom(3),
				isVegetarian: null.BoolFrom(true),
				isSet:        true,
			})).To(HaveLen(1))
			Expect(testDB.listRecipes(&filter{
				prepTimeTo: null.IntFrom(60),
				isSet:      true,
			})).To(HaveLen(4))
			Expect(testDB.listRecipes(&filter{
				prepTimeFrom: null.IntFrom(20),
				prepTimeTo:   null.IntFrom(60),
				isSet:        true,
			})).To(HaveLen(3))
			Expect(testDB.listRecipes(&filter{
				prepTimeFrom: null.IntFrom(20),
				prepTimeTo:   null.IntFrom(60),
				difficultyTo: null.IntFrom(3),
				isSet:        true,
			})).To(HaveLen(2))
			Expect(testDB.listRecipes(&filter{
				prepTimeFrom:   null.IntFrom(20),
				prepTimeTo:     null.IntFrom(60),
				difficultyFrom: null.IntFrom(2),
				difficultyTo:   null.IntFrom(4),
				isSet:          true,
			})).To(HaveLen(1))
			Expect(testDB.listRecipes(&filter{
				prepTimeFrom:   null.IntFrom(20),
				prepTimeTo:     null.IntFrom(60),
				difficultyFrom: null.IntFrom(2),
				difficultyTo:   null.IntFrom(4),
				isVegetarian:   null.BoolFrom(false),
				isSet:          true,
			})).To(HaveLen(0))
		})
	})
	Context("adding a new recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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
			testDB.sqlxDB.MustExec(`
			INSERT INTO hellofresh_user(hu_account, hu_access_token)
			VALUES
			('foo', 'faketoken')
			`)
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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
		It("adds a record in the recipe table and return the corresponding record", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			addedRecipe := testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name1"),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			Expect(addedRecipe.ID).To(Equal(1))
			Expect(addedRecipe.Name).To(Equal("name1"))
			Expect(addedRecipe.PrepareTime.Valid).To(BeFalse())
			Expect(addedRecipe.Difficulty.Valid).To(BeFalse())
			Expect(addedRecipe.IsVegetarian).To(BeFalse())
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(1))

			addedRecipe = testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name2"),
				PrepareTime:  null.IntFrom(2),
				Difficulty:   null.IntFrom(4),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			Expect(addedRecipe.ID).To(Equal(2))
			Expect(addedRecipe.Name).To(Equal("name2"))
			Expect(addedRecipe.PrepareTime.Int64).To(Equal(int64(2)))
			Expect(addedRecipe.Difficulty.Int64).To(Equal(int64(4)))
			Expect(addedRecipe.IsVegetarian).To(BeFalse())
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(2))
		})
		It("does nothing if the credential is not valid", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			actual := testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name1"),
				IsVegetarian: null.BoolFrom(false),
			}, "faild_token")
			Expect(actual).To(BeNil())
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(0))
		})
	})
	Context("updating a recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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

			testDB.sqlxDB.MustExec(`
			INSERT INTO hellofresh_user(hu_account, hu_access_token)
			VALUES
			('foo', 'faketoken')
			`)
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name1"),
				PrepareTime:  null.IntFrom(2),
				Difficulty:   null.IntFrom(4),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			actual := testDB.updateAndGetRecipeByCredential(&PutRecipeArg{
				Name:         null.StringFrom("name1_updated"),
				PrepareTime:  null.IntFrom(3),
				Difficulty:   null.IntFrom(4),
				IsVegetarian: null.BoolFrom(false),
			}, 1, "faketoken")

			Expect(actual.Name).To(Equal("name1_updated"))
			Expect(actual.PrepareTime.Int64).To(Equal(int64(3)))
			Expect(actual.Difficulty.Int64).To(Equal(int64(4)))
			Expect(actual.IsVegetarian).To(BeFalse())
		})
		It("does nothing if the recipe doesn't exist", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			actual := testDB.updateAndGetRecipeByCredential(&PutRecipeArg{
				Name:         null.StringFrom("name1_updated"),
				PrepareTime:  null.IntFrom(3),
				Difficulty:   null.IntFrom(4),
				IsVegetarian: null.BoolFrom(false),
			}, 2, "faketoken")

			Expect(actual).To(BeNil())
		})
		It("does nothing if the access to the recipe is not authorized", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			actual := testDB.updateAndGetRecipeByCredential(&PutRecipeArg{
				Name:         null.StringFrom("name1_updated"),
				PrepareTime:  null.IntFrom(3),
				Difficulty:   null.IntFrom(4),
				IsVegetarian: null.BoolFrom(false),
			}, 1, "failed_faketoken")

			Expect(actual).To(BeNil())
		})
	})
	Context("deleting a recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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
			testDB.sqlxDB.MustExec(`
			INSERT INTO hellofresh_user(hu_account, hu_access_token)
			VALUES
			('foo', 'faketoken')
			`)
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name1"),
				PrepareTime:  null.IntFrom(1),
				Difficulty:   null.IntFrom(2),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name2"),
				PrepareTime:  null.IntFrom(2),
				Difficulty:   null.IntFrom(4),
				IsVegetarian: null.BoolFrom(true),
			}, "faketoken")
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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
		It("deletes a existent record in the recipe table", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			deletedRecipe := testDB.deleteAndGetRecipeByCredential(1, "faketoken")
			Expect(testDB.getRecipeByID(1)).To(BeNil())
			Expect(testDB.getRecipeByID(2)).NotTo(BeNil())
			Expect(deletedRecipe.Name).To(Equal("name1"))
			Expect(deletedRecipe.PrepareTime.Int64).To(Equal(int64(1)))
			Expect(deletedRecipe.Difficulty.Int64).To(Equal(int64(2)))
			Expect(deletedRecipe.IsVegetarian).To(BeFalse())
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(1))
		})
		It("does nothing if the recipe doesn't exist", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			deletedRecipe := testDB.deleteAndGetRecipeByCredential(3, "faketoken")
			Expect(deletedRecipe).To(BeNil())
			Expect(testDB.getRecipeByID(1)).NotTo(BeNil())
			Expect(testDB.getRecipeByID(2)).NotTo(BeNil())
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(2))
		})
		It("does nothing if the access to the recipe is not authorized", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			deletedRecipe := testDB.deleteAndGetRecipeByCredential(1, "failed_token")
			Expect(deletedRecipe).To(BeNil())
			Expect(testDB.getRecipeByID(1)).NotTo(BeNil())
			Expect(testDB.getRecipeByID(2)).NotTo(BeNil())
			Expect(testDB.listRecipes(&filter{})).To(HaveLen(2))
		})
	})
	Context("rating a recipe", func() {
		BeforeEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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

			testDB.sqlxDB.MustExec(`
			INSERT INTO hellofresh_user(hu_account, hu_access_token)
			VALUES
			('foo', 'faketoken')
			`)
			testDB.addRecipeByCredential(&PostRecipeArg{
				Name:         null.StringFrom("name1"),
				PrepareTime:  null.IntFrom(2),
				Difficulty:   null.IntFrom(4),
				IsVegetarian: null.BoolFrom(false),
			}, "faketoken")
		})
		AfterEach(func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
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
		It("rates a existent recipe", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			actual := testDB.rateAndGetRecipe(&PostRateRecipeArg{
				Rating: null.IntFrom(3),
			}, 1)
			Expect(actual.Name).To(Equal("name1"))
			Expect(actual.RatedNum.Int64).To(Equal(int64(1)))
			Expect(actual.Rating.Float64).To(Equal(float64(3)))

			actual = testDB.rateAndGetRecipe(&PostRateRecipeArg{
				Rating: null.IntFrom(4),
			}, 1)
			Expect(actual.Name).To(Equal("name1"))
			Expect(actual.RatedNum.Int64).To(Equal(int64(2)))
			Expect(actual.Rating.Float64).To(Equal(float64(3.5)))

			actual = testDB.rateAndGetRecipe(&PostRateRecipeArg{
				Rating: null.IntFrom(5),
			}, 1)
			Expect(actual.Name).To(Equal("name1"))
			Expect(actual.RatedNum.Int64).To(Equal(int64(3)))
			Expect(actual.Rating.Float64).To(Equal(float64(4)))
		})
		It("does nothing if the recipe doesn't exist", func() {
			testDB := newSqlxPostgreSQL(testDBConnectionStringWithDatabase)
			defer testDB.close()

			actual := testDB.rateAndGetRecipe(&PostRateRecipeArg{
				Rating: null.IntFrom(3),
			}, 2)
			Expect(actual).To(BeNil())
		})
	})
})
