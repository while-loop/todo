package todo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/config"
	"github.com/while-loop/todo/pkg/tracker"
	"github.com/while-loop/todo/pkg/vcs"
)

type App struct {
	RepoMan    *vcs.Manager
	TrackerMan *tracker.Manager
	Router     *mux.Router
	Config     *config.Config
}

func New(config *config.Config, router *mux.Router) *App {
	rp := vcs.NewManager(config.VcsConfig)
	rp.ApplyRouter(router)

	return &App{
		RepoMan:    rp,
		TrackerMan: tracker.NewManager(config.TrackerConfig, rp.IssueChan()),
		Router:     router,
		Config:     config,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}

func (a *App) OnPush() error {
	return nil
}
