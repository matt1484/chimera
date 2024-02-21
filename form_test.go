package chimera_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/matt1484/chimera"
	"github.com/stretchr/testify/assert"
)

func TestAPIFormRequestReading(t *testing.T) {
	type SubStruct struct {
		Float float64 `form:"float"`
	}
	type Struct struct {
		Str   string    `form:"str"`
		Int   int       `form:"int"`
		Array []string  `form:"array"`
		Sub   SubStruct `form:"sub"`
	}
	body := Struct{
		Str:   "test",
		Int:   12345,
		Array: []string{"string1", "string2"},
		Sub: SubStruct{
			Float: -1.0,
		},
	}
	form := url.Values{
		"str":       []string{"test"},
		"int":       []string{"12345"},
		"array[0]":  []string{"string1"},
		"array[1]":  []string{"string2"},
		"sub.float": []string{"-1.0"},
	}

	api := chimera.NewAPI()
	primPath := addRequestTestHandler(t, api, http.MethodPost, testValidSimplePath+testValidLabelPath+testValidMatrixPath, &chimera.FormRequest[Struct, TestPrimitivePathParams]{Body: body, Params: testPrimitivePathParams})
	server := httptest.NewServer(api)
	resp, err := http.PostForm(server.URL+testValidPrimitiveSimplePathValues+testValidPrimitiveLabelPathValues+testValidPrimitiveMatrixPathValues, form)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*primPath).Params, testPrimitivePathParams)
	assert.Equal(t, (*primPath).Body, body)

	api = chimera.NewAPI()
	complexPath := addRequestTestHandler(t, api, http.MethodPost, testValidSimplePath+testValidLabelPath+testValidMatrixPath, &chimera.FormRequest[Struct, TestComplexPathParams]{Body: body, Params: testComplexPathParams})
	server = httptest.NewServer(api)
	resp, err = http.PostForm(server.URL+testValidComplexSimplePathValues+testValidComplexLabelPathValues+testValidComplexMatrixPathValues, form)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*complexPath).Params, testComplexPathParams)
	assert.Equal(t, (*complexPath).Body, body)

	api = chimera.NewAPI()
	primHeader := addRequestTestHandler(t, api, http.MethodPost, "/headertest", &chimera.FormRequest[Struct, TestPrimitiveHeaderParams]{Body: body, Params: testPrimitiveHeaderParams})
	server = httptest.NewServer(api)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/headertest", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
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
	complexHeader := addRequestTestHandler(t, api, http.MethodPost, "/headertest", &chimera.FormRequest[Struct, TestComplexHeaderParams]{Body: body, Params: testComplexHeaderParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/headertest", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	for k, v := range testValidSimpleComplexHeaderValues {
		req.Header.Set(k, v[0])
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*complexHeader).Params, testComplexHeaderParams)
	assert.Equal(t, (*complexHeader).Body, body)

	api = chimera.NewAPI()
	primCookie := addRequestTestHandler(t, api, http.MethodPost, "/cookietest", &chimera.FormRequest[Struct, TestPrimitiveCookieParams]{Body: body, Params: testPrimitiveCookieParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/cookietest", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	for _, c := range testValidFormPrimitiveCookieValues {
		req.AddCookie(&c)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*primCookie).Params, testPrimitiveCookieParams)
	assert.Equal(t, (*primCookie).Body, body)

	api = chimera.NewAPI()
	complexCookie := addRequestTestHandler(t, api, http.MethodPost, "/cookietest", &chimera.FormRequest[Struct, TestComplexCookieParams]{Body: body, Params: testComplexCookieParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/cookietest", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	for _, c := range testValidFormComplexCookieValues {
		req.AddCookie(&c)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*complexCookie).Params, testComplexCookieParams)
	assert.Equal(t, (*complexCookie).Body, body)

	api = chimera.NewAPI()
	primQuery := addRequestTestHandler(t, api, http.MethodPost, "/querytest", &chimera.FormRequest[Struct, TestPrimitiveQueryParams]{Body: body, Params: testPrimitiveQueryParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/querytest?"+testValidFormPrimitiveQueryValues.Encode(), strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*primQuery).Params, testPrimitiveQueryParams)
	assert.Equal(t, (*primQuery).Body, body)

	api = chimera.NewAPI()
	complexQuery := addRequestTestHandler(t, api, http.MethodPost, "/querytest", &chimera.FormRequest[Struct, TestComplexQueryParams]{Body: body, Params: testComplexQueryParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/querytest?"+testValidFormComplexQueryValues.Encode(), strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
	assert.Equal(t, (*complexQuery).Params, testComplexQueryParams)
	assert.Equal(t, (*complexQuery).Body, body)
}
