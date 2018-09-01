package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v3"
)

type mockDatastore struct {
	dataFunc func() interface{}
}

func (md *mockDatastore) listRecipes() []Recipe {
	return md.dataFunc().([]Recipe)
}

func (md *mockDatastore) addRecipe(arg PostRecipeArg) Recipe {
	return md.dataFunc().(Recipe)
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
		server := newTestAPIServer([]Recipe{
			{1, "name1", null.NewInt(0, false), null.NewInt(0, false), false},
			{11, "name11", null.NewInt(1, true), null.NewInt(2, true), true},
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
			   "is_vegetarian":false
			},
			{
			   "id":11,
			   "name":"name11",
			   "prepare_time":1,
			   "difficulty":2,
			   "is_vegetarian":true
			}
		]
		`))
	})

	It("lists empty results", func() {
		server := newTestAPIServer([]Recipe{})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/recipes", nil)
		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Get Recipes] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.MustArray()).To(HaveLen(0))
		Expect(rr.Body.String()).To(MatchJSON(`
		[]
		`))
	})
})

var _ = Describe("Adding a recipe", func() {
	It("adds a recipe and return the resulting JSON object", func() {
		server := newTestAPIServer(Recipe{32, "name3", null.NewInt(5, true), null.NewInt(0, false), false})
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

		server.httpServer.router.ServeHTTP(rr, req)

		jsonObj := newJSON(rr.Body.Bytes())
		GinkgoT().Logf("[Get Recipes] JSON Result: %s", jsonObj.pretty())
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonObj.Get("id").MustInt()).To(Equal(32))
		Expect(jsonObj.Get("name").MustString()).To(Equal("name3"))
		Expect(jsonObj.Get("prepare_time").MustInt()).To(Equal(5))
		Expect(jsonObj.Get("difficulty").Interface()).To(BeNil())
		Expect(jsonObj.Get("is_vegetarian").MustBool()).To(BeFalse())
	})

	It("responses with [400 Bad Request] when getting an invalid JSON argument", func() {
		server := newTestAPIServer(nil)
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

		server.httpServer.router.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusBadRequest))
	})
})
