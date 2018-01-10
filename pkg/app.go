package todo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/tracker"
	"github.com/while-loop/todo/pkg/vcs"
)

type App struct {
	RepoServices  []vcs.RepositoryService
	IssueTrackers []tracker.Tracker
	Router        *mux.Router
}

func New(repoServices []vcs.RepositoryService, issueTrackers []tracker.Tracker) *App {
	return &App{
		RepoServices:  repoServices,
		IssueTrackers: issueTrackers,
		Router:        mux.NewRouter(),
	}
}

func (a *App) Handler() http.Handler {
	return a.Router
}

func (a *App) OnPush() error {
	return nil
}
