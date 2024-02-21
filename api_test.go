package chimera_test

import (
	"net/http"

	"github.com/matt1484/chimera"
	"github.com/stretchr/testify/assert"
)

func addRequestTestHandler[ReqPtr chimera.RequestReaderPtr[Req], Req any](t assert.TestingT, api *chimera.API, method, path string, expectedReq ReqPtr) *ReqPtr {
	value := new(ReqPtr)
	switch method {
	case http.MethodGet:
		chimera.Get(api, path, func(req ReqPtr) (*chimera.EmptyResponse, error) {
			*value = req
			return nil, nil
		})
	case http.MethodPut:
		chimera.Put(api, path, func(req ReqPtr) (*chimera.EmptyResponse, error) {
			*value = req
			return nil, nil
		})
	case http.MethodPost:
		chimera.Post(api, path, func(req ReqPtr) (*chimera.EmptyResponse, error) {
			*value = req
			return nil, nil
		})
	case http.MethodDelete:
		chimera.Delete(api, path, func(req ReqPtr) (*chimera.EmptyResponse, error) {
			*value = req
			return nil, nil
		})
	case http.MethodPatch:
		chimera.Patch(api, path, func(req ReqPtr) (*chimera.EmptyResponse, error) {
			*value = req
			return nil, nil
		})
	}
	return value
}

func addResponseTestHandler[RespPtr chimera.ResponseWriterPtr[Resp], Resp any](t assert.TestingT, api *chimera.API, method, path string, resp RespPtr) {
	switch method {
	case http.MethodGet:
		chimera.Get(api, path, func(*chimera.EmptyRequest) (RespPtr, error) {
			return resp, nil
		})
	case http.MethodPut:
		chimera.Put(api, path, func(*chimera.EmptyRequest) (RespPtr, error) {
			return resp, nil
		})
	case http.MethodPost:
		chimera.Post(api, path, func(*chimera.EmptyRequest) (RespPtr, error) {
			return resp, nil
		})
	case http.MethodDelete:
		chimera.Delete(api, path, func(*chimera.EmptyRequest) (RespPtr, error) {
			return resp, nil
		})
	case http.MethodPatch:
		chimera.Patch(api, path, func(*chimera.EmptyRequest) (RespPtr, error) {
			return resp, nil
		})
	}
}

type TestValidCustomParam string

type TestStructParams struct {
	StringProp string `prop:"stringprop"`
	IntProp    int    `prop:"intprop"`
}
