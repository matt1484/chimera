package chimera_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matt1484/chimera"
	"github.com/stretchr/testify/assert"
)

func TestNoBodyResponseValid(t *testing.T) {
	api := chimera.NewAPI()
	addResponseTestHandler(t, api, http.MethodGet, "/headertest", &chimera.NoBodyResponse[TestPrimitiveHeaderParams]{Params: testPrimitiveHeaderParams})
	server := httptest.NewServer(api)
	resp, err := http.Get(server.URL + "/headertest")
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	for k, v := range testValidSimplePrimitiveHeaderValues {
		assert.Equal(t, resp.Header.Values(k)[0], v[0])
	}

	api = chimera.NewAPI()
	addResponseTestHandler(t, api, http.MethodGet, "/cookietest", &chimera.NoBodyResponse[TestPrimitiveCookieParams]{Params: testPrimitiveCookieParams})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + "/cookietest")
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	cookie := make(map[string]string)
	for _, c := range testValidFormPrimitiveCookieValues {
		cookie[c.Name] = c.Value
	}
	found := make(map[string]struct{})
	for _, c := range resp.Cookies() {
		assert.Equal(t, cookie[c.Name], c.Value)
		found[c.Name] = struct{}{}
	}
	assert.Equal(t, len(cookie), len(found))
}
