package chimera_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matt1484/chimera"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	api := chimera.NewAPI()
	called := make([]int, 0)
	api.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
		called = append(called, 1)
		return next(req)
	})
	api.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
		called = append(called, 2)
		return next(req)
	})
	group := api.Group("/base")
	group.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
		called = append(called, 3)
		return next(req)
	})

	sub := chimera.NewAPI()
	sub.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
		called = append(called, 4)
		return next(req)
	})
	chimera.Get(sub, "/route", func(*chimera.EmptyRequest) (*chimera.EmptyResponse, error) {
		return nil, nil
	})

	group.Mount("/sub", sub)
	server := httptest.NewServer(api)
	resp, _ := http.Get(server.URL + "/base/sub/route")
	server.Close()
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, called, []int{1, 2, 3, 4})
}
