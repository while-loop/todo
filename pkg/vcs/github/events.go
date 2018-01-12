package github

import (
	"net/http"
	"github.com/while-loop/todo/pkg/log"
	"encoding/json"
	"github.com/google/go-github/github"
	"fmt"
)

func (s *Service) handlePush(w http.ResponseWriter, r *http.Request) {
	var event github.PushEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Error("failed to decode github push event", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("got handlePush: %#v", event)
	// get changed files
	// foreach file, get all todos
	// send todos to issuechan (tracker will handle reducing and filtering)
}

func (s *Service) handleIssue(w http.ResponseWriter, r *http.Request) {
	var event github.IssueEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Error("failed to decode github push event", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("got handleIssue: %#v", event)
}

func (s *Service) handleInstallation(w http.ResponseWriter, r *http.Request) {
	var event github.InstallationEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Error("failed to decode github push event", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("got handleInstallation: %#v", event)
}
