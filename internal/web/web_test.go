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
	r := web.GetRouter(nil)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusMovedPermanently, resp.Code)
}

func TestHomeSlug(t *testing.T) {
	t.Skip()
	t.Parallel()
	r := web.GetRouter(nil)
	req, _ := http.NewRequest("GET", "/p/home", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
