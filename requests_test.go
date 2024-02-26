package chimera_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matt1484/chimera"
	"github.com/stretchr/testify/assert"
)

func TestNoBodyRequestValid(t *testing.T) {
	api := chimera.NewAPI()
	primPath := addRequestTestHandler(t, api, http.MethodGet, testValidSimplePath+testValidLabelPath+testValidMatrixPath, &chimera.NoBodyRequest[TestPrimitivePathParams]{Params: testPrimitivePathParams})
	server := httptest.NewServer(api)
	resp, err := http.Get(server.URL + testValidPrimitiveSimplePathValues + testValidPrimitiveLabelPathValues + testValidPrimitiveMatrixPathValues)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*primPath).Params, testPrimitivePathParams)

	api = chimera.NewAPI()
	complexPath := addRequestTestHandler(t, api, http.MethodGet, testValidSimplePath+testValidLabelPath+testValidMatrixPath, &chimera.NoBodyRequest[TestComplexPathParams]{testComplexPathParams})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + testValidComplexSimplePathValues + testValidComplexLabelPathValues + testValidComplexMatrixPathValues)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*complexPath).Params, testComplexPathParams)

	api = chimera.NewAPI()
	primHeader := addRequestTestHandler(t, api, http.MethodGet, "/headertest", &chimera.NoBodyRequest[TestPrimitiveHeaderParams]{Params: testPrimitiveHeaderParams})
	server = httptest.NewServer(api)
	req, err := http.NewRequest(http.MethodGet, server.URL+"/headertest", bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	for k, v := range testValidSimplePrimitiveHeaderValues {
		req.Header.Set(k, v[0])
	}
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*primHeader).Params, testPrimitiveHeaderParams)

	api = chimera.NewAPI()
	complexHeader := addRequestTestHandler(t, api, http.MethodGet, "/headertest", &chimera.NoBodyRequest[TestComplexHeaderParams]{Params: testComplexHeaderParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/headertest", bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	for k, v := range testValidSimpleComplexHeaderValues {
		req.Header.Set(k, v[0])
	}
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*complexHeader).Params, testComplexHeaderParams)

	api = chimera.NewAPI()
	primCookie := addRequestTestHandler(t, api, http.MethodGet, "/cookietest", &chimera.NoBodyRequest[TestPrimitiveCookieParams]{Params: testPrimitiveCookieParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/cookietest", bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	for _, c := range testValidFormPrimitiveCookieValues {
		req.AddCookie(&c)
	}
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*primCookie).Params, testPrimitiveCookieParams)

	api = chimera.NewAPI()
	complexCookie := addRequestTestHandler(t, api, http.MethodGet, "/cookietest", &chimera.NoBodyRequest[TestComplexCookieParams]{Params: testComplexCookieParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/cookietest", bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	for _, c := range testValidFormComplexCookieValues {
		req.AddCookie(&c)
	}
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*complexCookie).Params, testComplexCookieParams)

	api = chimera.NewAPI()
	primQuery := addRequestTestHandler(t, api, http.MethodGet, "/querytest", &chimera.NoBodyRequest[TestPrimitiveQueryParams]{Params: testPrimitiveQueryParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/querytest?"+testValidFormPrimitiveQueryValues.Encode(), bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*primQuery).Params, testPrimitiveQueryParams)

	api = chimera.NewAPI()
	complexQuery := addRequestTestHandler(t, api, http.MethodGet, "/querytest", &chimera.NoBodyRequest[TestComplexQueryParams]{Params: testComplexQueryParams})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/querytest?"+testValidFormComplexQueryValues.Encode(), bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*complexQuery).Params, testComplexQueryParams)
}

func TestInvalidParams(t *testing.T) {
	type PrimitiveInvalidPath struct {
		Param bool `param:"req,in=path"`
	}

	api := chimera.NewAPI()
	primPath := addRequestTestHandler(t, api, http.MethodGet, "/{req}", &chimera.NoBodyRequest[PrimitiveInvalidPath]{})
	server := httptest.NewServer(api)
	resp, err := http.Get(server.URL + "/test")
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *primPath)

	type SliceInvalidPath struct {
		Param []bool `param:"req,in=path"`
	}

	api = chimera.NewAPI()
	slicePath := addRequestTestHandler(t, api, http.MethodGet, "/{req}", &chimera.NoBodyRequest[SliceInvalidPath]{})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + "/test")
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *slicePath)

	type Struct struct {
		Param bool `prop:"x"`
	}

	type StructInvalidPath struct {
		Param Struct `param:"req,in=path"`
	}

	api = chimera.NewAPI()
	structPath := addRequestTestHandler(t, api, http.MethodGet, "/{req}", &chimera.NoBodyRequest[StructInvalidPath]{})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + "/x,test")
	server.Close()
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *structPath)

	type PrimitiveInvalidHeader struct {
		Param bool `param:"req,in=header,required"`
	}

	api = chimera.NewAPI()
	primHeader := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[PrimitiveInvalidHeader]{})
	server = httptest.NewServer(api)
	req, err := http.NewRequest(http.MethodGet, server.URL+"/test", bytes.NewBufferString(""))
	assert.NoError(t, err)
	req.Header.Add("req", "test")
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *primHeader)

	api = chimera.NewAPI()
	primHeader = addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[PrimitiveInvalidHeader]{})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/test", bytes.NewBufferString(""))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *primHeader)

	type SliceInvalidHeader struct {
		Param []bool `param:"req,in=header,required"`
	}

	api = chimera.NewAPI()
	sliceHeader := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[SliceInvalidHeader]{})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/test", bytes.NewBufferString(""))
	req.Header.Add("req", "test")
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *sliceHeader)

	type StructInvalidHeader struct {
		Param Struct `param:"req,in=header"`
	}

	api = chimera.NewAPI()
	structHeader := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[StructInvalidHeader]{})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/test", bytes.NewBufferString(""))
	req.Header.Add("req", "x,test")
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *structHeader)

	type PrimitiveInvalidCookie struct {
		Param bool `param:"req,in=cookie,required"`
	}

	api = chimera.NewAPI()
	primCookie := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[PrimitiveInvalidCookie]{})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/test", bytes.NewBufferString(""))
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "req",
		Value: "test",
	})
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *primCookie)

	api = chimera.NewAPI()
	primCookie = addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[PrimitiveInvalidCookie]{})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/test", bytes.NewBufferString(""))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *primCookie)

	type SliceInvalidCookie struct {
		Param []bool `param:"req,in=cookie,required"`
	}

	api = chimera.NewAPI()
	sliceCookie := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[SliceInvalidCookie]{})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/test", bytes.NewBufferString(""))
	req.AddCookie(&http.Cookie{
		Name:  "req",
		Value: "test",
	})
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *sliceCookie)

	type StructInvalidCookie struct {
		Param Struct `param:"req,in=cookie"`
	}

	api = chimera.NewAPI()
	structCookie := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[StructInvalidCookie]{})
	server = httptest.NewServer(api)
	req, err = http.NewRequest(http.MethodGet, server.URL+"/test", bytes.NewBufferString(""))
	req.AddCookie(&http.Cookie{
		Name:  "req",
		Value: "x,test",
	})
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *structCookie)

	type PrimitiveInvalidQuery struct {
		Param bool `param:"req,in=query,required"`
	}

	api = chimera.NewAPI()
	primQuery := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[PrimitiveInvalidQuery]{})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + "/test?req=test")
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *primQuery)

	api = chimera.NewAPI()
	primQuery = addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[PrimitiveInvalidQuery]{})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + "/test")
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *primQuery)

	type SliceInvalidQuery struct {
		Param []bool `param:"req,in=query,required"`
	}

	api = chimera.NewAPI()
	sliceQuery := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[SliceInvalidQuery]{})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + "/test?req=test")
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *sliceQuery)

	type StructInvalidQuery struct {
		Param Struct `param:"req,in=query"`
	}

	api = chimera.NewAPI()
	structQuery := addRequestTestHandler(t, api, http.MethodGet, "/test", &chimera.NoBodyRequest[StructInvalidQuery]{})
	server = httptest.NewServer(api)
	resp, err = http.Get(server.URL + "/test?req=x,test")
	assert.NoError(t, err)
	server.Close()
	assert.Equal(t, resp.StatusCode, 422)
	assert.Nil(t, *structQuery)
}

func TestOpenAPI(t *testing.T) {
	api := chimera.NewAPI()
	route := chimera.Get(api, "/testopenapi", func(*chimera.EmptyRequest) (*chimera.EmptyResponse, error) {
		return nil, nil
	})
	assert.Contains(t, api.OpenAPISpec().Paths, "/testopenapi")
	route.Internalize()
	assert.NotContains(t, api.OpenAPISpec().Paths, "/testopenapi")
}
