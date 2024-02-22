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
	firstCalled := false
	secondCalled := false
	thirdCalled := false
	api.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
		firstCalled = true
		return next(req)
	})
	api.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
		secondCalled = true
		return next(req)
	})
	group := api.Group("/base")
	chimera.Get(group, "/sub", func(*chimera.EmptyRequest) (*chimera.EmptyResponse, error) {
		return nil, nil
	})
	group.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
		thirdCalled = true
		return next(req)
	})

	server := httptest.NewServer(api)
	resp, _ := http.Get(server.URL + "/base/sub")
	server.Close()
	assert.Equal(t, resp.StatusCode, 200)
	assert.True(t, firstCalled)
	assert.True(t, secondCalled)
	assert.True(t, thirdCalled)
}
