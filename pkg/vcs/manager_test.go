package vcs

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/while-loop/todo/pkg/vcs/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testVCS struct{}

func (t *testVCS) Name() string {
	return "testvcs"
}

func (t *testVCS) Init(webhookRouter *mux.Router) {
	webhookRouter.HandleFunc("/testvcs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello test")
	}).Methods(http.MethodPost)
}

func TestWebHookSubRouter(t *testing.T) {
	router := mux.NewRouter()
	man := NewManager(&config.VcsConfig{
		Github: &config.GithubConfig{},
	}, router)
	man.services["testvcs"] = &testVCS{}
	man.initRouter(router)

	r := httptest.NewRequest(http.MethodPost, "/webhook/testvcs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello test", string(body))

	r = httptest.NewRequest(http.MethodPost, "/webhook/github", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "", string(body))
}
