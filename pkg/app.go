package todo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/config"
	"github.com/while-loop/todo/pkg/tracker"
	"github.com/while-loop/todo/pkg/vcs"
	"github.com/while-loop/todo/pkg/log"
)

type App struct {
	RepoMan    *vcs.Manager
	TrackerMan *tracker.Manager
	Router     *mux.Router
	Config     *config.Config
}

func New(config *config.Config, router *mux.Router, logger log.AnalysisLogger) *App {
	tm := tracker.NewManager(config.TrackerConfig, logger)
	rp := vcs.NewManager(config.VcsConfig, router, tm, logger)
	return &App{
		RepoMan:    rp,
		TrackerMan: tm,
		Router:     router,
		Config:     config,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}
