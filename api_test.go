package main

import (
	"net/http"
	"net/http/httptest"

	json "github.com/bitly/go-simplejson"
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

		jsonBody, _ := json.NewJson(rr.Body.Bytes())
		prettyJSON, _ := jsonBody.EncodePretty()
		GinkgoT().Logf("[Get Recipes] JSON Result: %s", prettyJSON)
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonBody.MustArray()).To(HaveLen(2))
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

		jsonBody, _ := json.NewJson(rr.Body.Bytes())
		prettyJSON, _ := jsonBody.EncodePretty()
		GinkgoT().Logf("[Get Recipes] JSON Result: %s", prettyJSON)
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(jsonBody.MustArray()).To(HaveLen(0))
		Expect(rr.Body.String()).To(MatchJSON(`
		[]
		`))
	})
})
