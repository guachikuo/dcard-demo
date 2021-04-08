package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

var (
	mockIP = "192.168.1.1"
)

func mockHTTPRequest(method string, path string, reqBody interface{}) *http.Request {
	reqBodyReader := bytes.NewBufferString("nilBody")
	if reqBody != nil {
		reqBodyReader = new(bytes.Buffer)
		json.NewEncoder(reqBodyReader).Encode(reqBody)
	}
	req, _ := http.NewRequest(method, path, reqBodyReader)
	req.Header.Add("Content-Type", "application/json; charset=uft-8")
	req.Header.Add("X-Forwarded-For", mockIP)
	return req
}

func newRecorder() *httptest.ResponseRecorder {
	respRecorder := httptest.NewRecorder()
	respRecorder.Code = -1

	return respRecorder
}
