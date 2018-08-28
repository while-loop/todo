package github

import (
	"context"
	"net/http"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io/ioutil"
	"strings"

	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/parser"
	"github.com/while-loop/todo/pkg/vcs/config"
)

const (
	name                      = "github"
	issues                    = "issues"
	push                      = "push"
	installation              = "installation"
	installation_repositories = "installation_repositories"
)

type Service struct {
	router        *mux.Router
	eventHandlers map[string]http.HandlerFunc
	config        *config.GithubConfig
	parser        parser.TodoParser
	issueCreator  issue.Creator
	logger        log.AnalysisLogger
}

func NewService(config *config.GithubConfig, creator issue.Creator, logger log.AnalysisLogger) *Service {
	s := &Service{
		issueCreator:  creator,
		router:        mux.NewRouter(),
		eventHandlers: map[string]http.HandlerFunc{},
		config:        config,
		parser:        parser.New(), // todo(while-loop) add parser as param
		logger:        logger,
	}

	s.eventHandlers[issues] = s.handleIssue
	s.eventHandlers[installation] = s.handleInstallation
	s.eventHandlers[installation_repositories] = s.handleRepoInstallation
	s.eventHandlers[push] = s.handlePush
	return s
}

func (s *Service) Name() string {
	return name
}

func (s *Service) Init(webhookRouter *mux.Router) {
	webhookRouter.HandleFunc("/"+name, s.webhook).Methods("POST")
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
