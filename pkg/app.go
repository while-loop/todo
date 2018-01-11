package todo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/config"
	"github.com/while-loop/todo/pkg/publisher"
	"github.com/while-loop/todo/pkg/tracker"
	"github.com/while-loop/todo/pkg/vcs"
)

type App struct {
	RepoMan      *vcs.Manager
	TrackerMan   *tracker.Manager
	PublisherMan *publisher.Manager
	Router       *mux.Router
	Config       *config.Config
}

func New(config *config.Config) *App {
	router := mux.NewRouter()
	rp := vcs.NewManager(config.VcsConfig)
	rp.ApplyRouter(router)

	return &App{
		RepoMan:      rp,
		TrackerMan:   tracker.NewManager(config.TrackerConfig),
		PublisherMan: publisher.NewManager(config.PublisherConfig),
		Router:       router,
		Config:       config,
	}
}

func (a *App) Handler() http.Handler {
	return a.Router
}

func (a *App) OnPush() error {
	return nil
}
