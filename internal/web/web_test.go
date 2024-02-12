package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bnadland/rhizome/internal/web"
	"github.com/stretchr/testify/assert"
)

func TestHomepage(t *testing.T) {
	t.Parallel()
	r := web.NewRouter(nil)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusFound, resp.Code)
}
