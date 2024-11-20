package test

import (
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var routerHandler http.Handler

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	if routerHandler == nil {
		routerHandler = k.Router.Handler()
	}

	routerHandler.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func checkResponse(t *testing.T, expectedStatusCode int, expectedResponse string, response *httptest.ResponseRecorder, formats ...interface{}) {
	ja := jsonassert.New(t)
	checkResponseCode(t, expectedStatusCode, response.Code)

	receivedResponse := response.Body.String()

	if receivedResponse == "" {
		assert.Equal(t, expectedResponse, receivedResponse)
		return
	}

	if formats != nil {
		ja.Assertf(receivedResponse, expectedResponse, formats)
	} else {
		ja.Assertf(receivedResponse, expectedResponse)
	}
}
