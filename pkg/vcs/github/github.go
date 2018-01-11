package github

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/tracker"
	"github.com/while-loop/todo/pkg/vcs/config"
	"golang.org/x/oauth2"
)

const (
	name = "github"
)

type Service struct {
	router   *mux.Router
	ghClient *github.Client
	issueCh  <-chan tracker.Issue
}

func NewService(config *config.GithubConfig, issueChan <-chan tracker.Issue) *Service {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.AccessToken})
	oauthClient := oauth2.NewClient(context.Background(), ts)
	s := &Service{
		issueCh:  issueChan,
		router:   mux.NewRouter(),
		ghClient: github.NewClient(oauthClient),
	}

	return s
}

func (s *Service) GetRepository(user, project string) error {
	panic("implement me")
}

func (s *Service) Name() string {
	return name
}

func (s *Service) Handler() http.Handler {
	s.router.HandleFunc("/webhook/"+name, s.webhook)
	return s.router
}

func (s *Service) webhook(w http.ResponseWriter, r *http.Request) {
	log.Info(name, r.URL, w, r)
}
