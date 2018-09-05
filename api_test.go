package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	null "gopkg.in/guregu/null.v3"
)

type mockDatastore struct {
	dataFunc func() interface{}
}

func (md *mockDatastore) listRecipes(f *filter) []*Recipe {
	return md.dataFunc().([]*Recipe)
}

func (md *mockDatastore) addRecipeByCredential(arg *PostRecipeArg, token string) *Recipe {
	if d := md.dataFunc(); d != nil {
		return md.dataFunc().(*Recipe)
	}
	return nil
}

func (md *mockDatastore) getRecipeByID(id int) *Recipe {
	if d := md.dataFunc(); d != nil {
		return md.dataFunc().(*Recipe)
	}
	return nil
}

func (md *mockDatastore) updateAndGetRecipeByCredential(arg *PutRecipeArg, id int, token string) *Recipe {
	if d := md.dataFunc(); d != nil {
		return md.dataFunc().(*Recipe)
	}
	return nil
}

func (md *mockDatastore) deleteAndGetRecipeByCredential(id int, token string) *Recipe {
	if d := md.dataFunc(); d != nil {
		return md.dataFunc().(*Recipe)
	}
	return nil
}

func (md *mockDatastore) rateAndGetRecipe(arg *PostRateRecipeArg, id int) *Recipe {
	if d := md.dataFunc(); d != nil {
		return md.dataFunc().(*Recipe)
	}
	return nil
}

func newTestAPIServer(data interface{}) *apiServer {
	md := &mockDatastore{
		dataFunc: func() interface{} {
			return data
		},
	}
	s := &apiServer{
		httpServer: newGinHTTPServer(),
		datastore:  md,
	}
	s.routes()
	return s
}

var _ = Describe("Listing recipes", func() {
	It("lists non-empty results", func() {
		server := newTestAPIServer([]*Recipe{
			{1, "name1", null.IntFromPtr(nil), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)},
			{11, "name11", null.IntFrom(1), null.IntFrom(2), true, null.FloatFrom(0.0), null.IntFrom(0)},
		})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/recipes", nil)
		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Get Recipes] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.MustArray()).To(HaveLen(2))
		Expect(rr.Body.String()).To(MatchJSON(`
		[
			{
			   "id":1,
			   "name":"name1",
			   "prepare_time":null,
			   "difficulty":null,
			   "is_vegetarian":false,
			   "rating": 0,
			   "rated_num": 0
			},
			{
			   "id":11,
			   "name":"name11",
			   "prepare_time":1,
			   "difficulty":2,
			   "is_vegetarian":true,
			   "rating": 0,
			   "rated_num": 0
			}
		]
		`))
	})
	It("lists empty results", func() {
		server := newTestAPIServer([]*Recipe{})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/recipes", nil)
		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Get Recipes] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.MustArray()).NotTo(BeNil())
		Expect(jsonObj.MustArray()).To(HaveLen(0))
		Expect(rr.Body.String()).To(MatchJSON(`
		[]
		`))
	})
})

var _ = Describe("Adding a recipe", func() {
	It("adds a recipe and returns the resulting JSON object", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/recipes", newJSON([]byte(`
		{
			"name":"name3",
			"prepare_time":5,
			"difficulty":null,
			"is_vegetarian":false
		}
		`)).buffer())
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Add A Recipe] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.Get("id").MustInt()).To(Equal(32))
		Expect(jsonObj.Get("name").MustString()).To(Equal("name3"))
		Expect(jsonObj.Get("prepare_time").MustInt()).To(Equal(5))
		Expect(jsonObj.Get("difficulty").Interface()).To(BeNil())
		Expect(jsonObj.Get("is_vegetarian").MustBool()).To(BeFalse())
	})
	It("responses with [404 Not Found] when the user's credential is not valid", func() {
		server := newTestAPIServer(nil)
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer([]byte(`
		{
			"name":"name3",
			"prepare_time":5,
			"difficulty":null,
			"is_vegetarian":false
		}
		`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
	It("responses with [400 Bad Request] when getting an invalid JSON argument", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer([]byte(`
		{
			"name":"name3",
			"prepare_time":5,
			"difficulty":null,
			"is_vegetarian":false,
		}
		`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusBadRequest))
	})
})

var _ = Describe("Getting a recipe by ID", func() {
	It("gets a recipe and returns the corresponding JSON object", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/recipes/32", nil)

		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Get A Recipe By ID] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.Get("id").MustInt()).To(Equal(32))
		Expect(jsonObj.Get("name").MustString()).To(Equal("name3"))
		Expect(jsonObj.Get("prepare_time").MustInt()).To(Equal(5))
		Expect(jsonObj.Get("difficulty").Interface()).To(BeNil())
		Expect(jsonObj.Get("is_vegetarian").MustBool()).To(BeFalse())
	})
	It("responses with [404 Not Found] when getting an invalid parameter", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/recipes/ff", nil)

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
	It("responses with [404 Not Found] when the recipe doesn't exist", func() {
		server := newTestAPIServer(nil)
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/recipes/32", nil)

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
})

var _ = Describe("Updating a recipe by ID", func() {
	It("updates a recipe and gets the updated JSON object", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFrom(3), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipes/32", bytes.NewBuffer([]byte(`
		{
			"prepare_time":5,
			"difficulty":3
		}
		`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Update A Recipe By ID] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.Get("prepare_time").MustInt()).To(Equal(5))
		Expect(jsonObj.Get("difficulty").MustInt()).To(Equal(3))
	})
	It("responses with [404 Not Found] when getting an invalid parameter", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipes/ff", bytes.NewBuffer([]byte(`
		{
			"prepare_time":5,
			"difficulty":3
		}
		`)))
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
	It("responses with [404 Not Found] when the recipe is not authorized or not found", func() {
		server := newTestAPIServer(nil)
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipes/32", bytes.NewBuffer([]byte(`
		{
			"prepare_time":5,
			"difficulty":3
		}
		`)))
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
	It("responses with [400 Bad Request] if the JSON argument is invalid", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipes/32", bytes.NewBuffer([]byte(`
		{
			"name":"name3",
			"prepare_time":5,
			"difficulty":null,
			"is_vegetarian":false,
		}
		`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusBadRequest))
	})
})

var _ = Describe("Deleting a recipe by ID", func() {
	It("deletes a recipe and gets the deleted JSON object", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/recipes/32", nil)
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Delete A Recipe By ID] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.Get("id").MustInt()).To(Equal(32))
		Expect(jsonObj.Get("name").MustString()).To(Equal("name3"))
		Expect(jsonObj.Get("prepare_time").MustInt()).To(Equal(5))
		Expect(jsonObj.Get("difficulty").Interface()).To(BeNil())
		Expect(jsonObj.Get("is_vegetarian").MustBool()).To(BeFalse())
	})
	It("responses with [404 Not Found] when getting an invalid parameter", func() {
		server := newTestAPIServer(&Recipe{32, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(0.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/recipes/ff", nil)
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
	It("responses with [404 Not Found] when the recipe is not authorized or not found", func() {
		server := newTestAPIServer(nil)
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/recipes/32", nil)
		req.Header.Set("Authorization", "faketoken")

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
})

var _ = Describe("Rating a recipe by ID", func() {
	It("Rates a recipe and gets the updated JSON object", func() {
		server := newTestAPIServer(&Recipe{3, "name3", null.IntFrom(5), null.IntFrom(3), false, null.FloatFrom(3.0), null.IntFrom(1)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipes/3", bytes.NewBuffer([]byte(`
		{
			"rating":3
		}
		`)))
		req.Header.Set("Content-Type", "application/json")

		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Rate A Recipe By ID] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.Get("rating").MustFloat64()).To(Equal(3.0))
		Expect(jsonObj.Get("rated_num").MustInt()).To(Equal(1))
	})
	It("responses with [404 Not Found] when getting an invalid parameter", func() {
		server := newTestAPIServer(&Recipe{3, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(3.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipes/ff", bytes.NewBuffer([]byte(`
		{
			"rating":3
		}
		`)))

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
	It("responses with [404 Not Found] when the recipe is not found", func() {
		server := newTestAPIServer(nil)
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipes/3", bytes.NewBuffer([]byte(`
		{
			"rating":3
		}
		`)))

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})
	It("responses with [400 Bad Request] if the JSON argument is invalid", func() {
		server := newTestAPIServer(&Recipe{3, "name3", null.IntFrom(5), null.IntFromPtr(nil), false, null.FloatFrom(3.0), null.IntFrom(0)})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipes/3", bytes.NewBuffer([]byte(`
		{
			"rating":3,
		}
		`)))
		req.Header.Set("Content-Type", "application/json")

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusBadRequest))
	})
})
