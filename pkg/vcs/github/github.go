package github

import (
	"context"
	"net/http"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/parser"
	"github.com/while-loop/todo/pkg/tracker"
	"github.com/while-loop/todo/pkg/vcs/config"
	"golang.org/x/oauth2"
	"hash"
	"io/ioutil"
	"strings"
)

const (
	name         = "github"
	issues       = "issues"
	push         = "push"
	installation = "installation"
)

type Service struct {
	router        *mux.Router
	ghClient      *github.Client
	issueCh       chan<- tracker.Issue
	eventHandlers map[string]http.HandlerFunc
	config        *config.GithubConfig
	parser        parser.TodoParser
}

func NewService(config *config.GithubConfig, issueChan chan<- tracker.Issue) *Service {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.AccessToken})
	oauthClient := oauth2.NewClient(context.Background(), ts)
	s := &Service{
		issueCh:       issueChan,
		router:        mux.NewRouter(),
		ghClient:      github.NewClient(oauthClient),
		eventHandlers: map[string]http.HandlerFunc{},
		config:        config,
		parser:        parser.New(), // todo(while-loop) add parser as param
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
	event := r.Header.Get("X-GitHub-Event")
	sig := r.Header.Get("X-Hub-Signature")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("failed to read github webhook body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !validBody(body, s.config.WebhookSecret, sig) {
		log.Error("failed to verify payload hash", sig)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if h, exists := s.eventHandlers[event]; exists {
		ctx := context.WithValue(r.Context(), "body", body)
		h(w, r.WithContext(ctx))
	} else {
		log.Errorf("Unknown event called: %s", event)
		w.WriteHeader(http.StatusNotFound)
	}
}

func validBody(body []byte, secret string, sig string) bool {
	split := strings.Split(sig, "=")
	if len(split) != 2 {
		log.Error("failed to get hash func and sig from github", sig)
		return false
	}

	var hAlg func() hash.Hash
	switch split[0] {
	case "sha1":
		fallthrough
	default:
		hAlg = sha1.New
	}

	h := hmac.New(hAlg, []byte(secret))
	n, err := h.Write(body)
	if n != len(body) || err != nil {
		log.Errorf("failed to write to hmac. %v sig: %s", err, sig)
		return false
	}

	bs, err := hex.DecodeString(split[1])
	if err != nil {
		log.Error("failed to decode github hash ", split[1])
		return false
	}
	return hmac.Equal(h.Sum(nil), bs)
}
