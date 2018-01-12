package github

import (
	"context"
	"net/http"

	"io/ioutil"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/tracker"
	"github.com/while-loop/todo/pkg/vcs/config"
	"golang.org/x/oauth2"
)

const (
	name = "github"
	issues = "issues"
	push = "push"
	installation = "installation"
)

type Service struct {
	router        *mux.Router
	ghClient      *github.Client
	issueCh       <-chan tracker.Issue
	eventHandlers map[string]http.HandlerFunc
}

func NewService(config *config.GithubConfig, issueChan <-chan tracker.Issue) *Service {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.AccessToken})
	oauthClient := oauth2.NewClient(context.Background(), ts)
	s := &Service{
		issueCh:       issueChan,
		router:        mux.NewRouter(),
		ghClient:      github.NewClient(oauthClient),
		eventHandlers: map[string]http.HandlerFunc{},
	}

	s.eventHandlers[issues] = s.handleIssue
	s.eventHandlers[installation] = s.handleInstallation
	s.eventHandlers[push] = s.handlePush
	return s
}

func (s *Service) GetRepository(user, project string) error {
	panic("implement me")
}

func (s *Service) Name() string {
	return name
}

func (s *Service) Handler() http.Handler {
	s.router.HandleFunc("/webhook/"+name, s.webhook).Methods("POST")
	return s.router
}

func (s *Service) webhook(w http.ResponseWriter, r *http.Request) {
	log.Info(name, r.URL, r.Header)
	bs, err := ioutil.ReadAll(r.Body)
	log.Info(err, "\n", string(bs))
	event := r.Header.Get("X-GitHub-Event")

	if h, exists := s.eventHandlers[event]; exists {
		h(w, r)
	} else {
		log.Errorf("Unknown event called: %s", event)
		w.WriteHeader(http.StatusNotFound)
	}
}
