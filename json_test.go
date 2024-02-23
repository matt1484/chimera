package chimera_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matt1484/chimera"
	"github.com/stretchr/testify/assert"
)

func TestJSONRequestValid(t *testing.T) {
	type Struct struct {
		Str string `json:"str"`
		Int int    `json:"int"`
	}
	body := Struct{
		Str: "a test",
		Int: 12345,
	}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	api := chimera.NewAPI()
	primPath := addRequestTestHandler(t, api, http.MethodPost, testValidSimplePath+testValidLabelPath+testValidMatrixPath, &chimera.JSON[Struct, TestPrimitivePathParams]{Body: body, Params: testPrimitivePathParams})
	server := httptest.NewServer(api)
	resp, err := http.Post(server.URL+testValidPrimitiveSimplePathValues+testValidPrimitiveLabelPathValues+testValidPrimitiveMatrixPathValues, "application/json", bytes.NewBuffer(b))
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*primPath).Params, testPrimitivePathParams)
	assert.Equal(t, (*primPath).Body, body)

	api = chimera.NewAPI()
	complexPath := addRequestTestHandler(t, api, http.MethodPost, testValidSimplePath+testValidLabelPath+testValidMatrixPath, &chimera.JSONRequest[Struct, TestComplexPathParams]{Body: body, Params: testComplexPathParams})
	server = httptest.NewServer(api)
	resp, err = http.Post(server.URL+testValidComplexSimplePathValues+testValidComplexLabelPathValues+testValidComplexMatrixPathValues, "application/json", bytes.NewBuffer(b))
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*complexPath).Params, testComplexPathParams)
	assert.Equal(t, (*complexPath).Body, body)

	api = chimera.NewAPI()
	primHeader := addRequestTestHandler(t, api, http.MethodPost, "/headertest", &chimera.JSONRequest[Struct, TestPrimitiveHeaderParams]{Body: body, Params: testPrimitiveHeaderParams})
	server = httptest.NewServer(api)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/headertest", bytes.NewBuffer(b))
	assert.NoError(t, err)
	for k, v := range testValidSimplePrimitiveHeaderValues {
		req.Header.Set(k, v[0])
	}
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*primHeader).Params, testPrimitiveHeaderParams)
	assert.Equal(t, (*primHeader).Body, body)

	api = chimera.NewAPI()
	complexHeader := addRequestTestHandler(t, api, http.MethodPost, "/headertest", &chimera.JSONRequest[Struct, TestComplexHeaderParams]{Body: body, Params: testComplexHeaderParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/headertest", bytes.NewBuffer(b))
	assert.NoError(t, err)
	for k, v := range testValidSimpleComplexHeaderValues {
		req.Header.Set(k, v[0])
	}
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*complexHeader).Params, testComplexHeaderParams)
	assert.Equal(t, (*complexHeader).Body, body)

	api = chimera.NewAPI()
	primCookie := addRequestTestHandler(t, api, http.MethodPost, "/cookietest", &chimera.JSONRequest[Struct, TestPrimitiveCookieParams]{Body: body, Params: testPrimitiveCookieParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/cookietest", bytes.NewBuffer(b))
	assert.NoError(t, err)
	for _, c := range testValidFormPrimitiveCookieValues {
		req.AddCookie(&c)
	}
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*primCookie).Params, testPrimitiveCookieParams)
	assert.Equal(t, (*primCookie).Body, body)

	api = chimera.NewAPI()
	complexCookie := addRequestTestHandler(t, api, http.MethodPost, "/cookietest", &chimera.JSONRequest[Struct, TestComplexCookieParams]{Body: body, Params: testComplexCookieParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/cookietest", bytes.NewBuffer(b))
	assert.NoError(t, err)
	for _, c := range testValidFormComplexCookieValues {
		req.AddCookie(&c)
	}
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*complexCookie).Params, testComplexCookieParams)
	assert.Equal(t, (*complexCookie).Body, body)

	api = chimera.NewAPI()
	primQuery := addRequestTestHandler(t, api, http.MethodPost, "/querytest", &chimera.JSONRequest[Struct, TestPrimitiveQueryParams]{Body: body, Params: testPrimitiveQueryParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/querytest?"+testValidFormPrimitiveQueryValues.Encode(), bytes.NewBuffer(b))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*primQuery).Params, testPrimitiveQueryParams)
	assert.Equal(t, (*primQuery).Body, body)

	api = chimera.NewAPI()
	complexQuery := addRequestTestHandler(t, api, http.MethodPost, "/querytest", &chimera.JSONRequest[Struct, TestComplexQueryParams]{Body: body, Params: testComplexQueryParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/querytest?"+testValidFormComplexQueryValues.Encode(), bytes.NewBuffer(b))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*complexQuery).Params, testComplexQueryParams)
	assert.Equal(t, (*complexQuery).Body, body)
}

func TestJSONResponseValid(t *testing.T) {
	type Struct struct {
		Str string `json:"str"`
		Int int    `json:"int"`
	}
	body := Struct{
		Str: "a test",
		Int: 12345,
	}
	api := chimera.NewAPI()
	addResponseTestHandler(t, api, http.MethodGet, "/headertest", &chimera.JSON[Struct, TestPrimitiveHeaderParams]{Body: body, Params: testPrimitiveHeaderParams})
	server := httptest.NewServer(api)
	resp, err := http.Get(server.URL + "/headertest")
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	for k, v := range testValidSimplePrimitiveHeaderValues {
		assert.Equal(t, resp.Header.Values(k)[0], v[0])
	}
	server.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"str":"a test","int":12345}`, string(b))

	api = chimera.NewAPI()
	addResponseTestHandler(t, api, http.MethodGet, "/cookietest", &chimera.JSONResponse[Struct, TestPrimitiveCookieParams]{Body: body, Params: testPrimitiveCookieParams})
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
	assert.Equal(t, `{"str":"a test","int":12345}`, string(b))
}
