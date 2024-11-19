package test

import (
	kernel "github.com/mik3lon/starter-template/internal/pkg/infrastructure/kernel"
	"github.com/mik3lon/starter-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type HttpHeaders map[string]string

var k *kernel.Kernel

func setUp() {
	if k == nil {
		cnf := config.LoadTestConfig()
		k = kernel.Init(cnf)
	}

	wg := &sync.WaitGroup{}
	wg.Wait()
}

func executeJsonApiRequest(t *testing.T, method, url string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)

	if len(headers) != 0 {
		for headerName, value := range headers {
			req.Header.Set(headerName, value)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	return executeRequest(req)
}
