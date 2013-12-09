package github

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestGetRepo_returns_AHydratedRepo(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	}))
	defer fakeServer.Close()

	repo, header := GetRepos("http://someplace.com")

	assert.NotNil(t, repo)
	assert.NotNil(t, header)
}
