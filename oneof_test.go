package chimera_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matt1484/chimera"
	"github.com/stretchr/testify/assert"
)

type TestOneOfBodyX struct {
	S string `json:"s"`
}

type TestOneOfBodyY struct {
	B bool `json:"b"`
}

type TestOneOfValidStruct struct {
	X *chimera.JSONResponse[TestOneOfBodyX, chimera.Nil] `response:"statusCode=206"`
	Y *chimera.JSONResponse[TestOneOfBodyY, chimera.Nil] `response:"statusCode=207"`
}

func TestValidOneOfRequest(t *testing.T) {
	api := chimera.NewAPI()
	addResponseTestHandler(t, api, http.MethodGet, "/testoneof", &chimera.OneOfResponse[TestOneOfValidStruct]{
		Response: TestOneOfValidStruct{
			X: &chimera.JSONResponse[TestOneOfBodyX, chimera.Nil]{
				Body: TestOneOfBodyX{
					S: "test x",
				},
			},
		},
	})
	server := httptest.NewServer(api)
	resp, err := http.Get(server.URL + "/testoneof")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 206)
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, b, []byte(`{"s":"test x"}`))

	api = chimera.NewAPI()
	addResponseTestHandler(t, api, http.MethodGet, "/testoneof", &chimera.OneOfResponse[TestOneOfValidStruct]{
		Response: TestOneOfValidStruct{
			Y: &chimera.JSONResponse[TestOneOfBodyY, chimera.Nil]{
				Body: TestOneOfBodyY{
					B: true,
				},
			},
		},
	})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + "/testoneof")
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 207)
	b, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, b, []byte(`{"b":true}`))
}
