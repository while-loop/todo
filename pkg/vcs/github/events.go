package github

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"github.com/luci/go-render/render"
	"github.com/while-loop/todo/pkg/log"
	"net/http"
)

func (s *Service) handleIssue(w http.ResponseWriter, r *http.Request) {
	var event github.IssueEvent
	if err := json.Unmarshal(r.Context().Value("body").([]byte), &event); err != nil {
		log.Error("failed to decode github IssueEvent", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("got handleIssue: %s", render.Render(event))
}

func (s *Service) handleInstallation(w http.ResponseWriter, r *http.Request) {
	var event github.InstallationEvent
	if err := json.Unmarshal(r.Context().Value("body").([]byte), &event); err != nil {
		log.Error("failed to decode github InstallationEvent", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("got handleInstallation: %s", render.Render(event))
}
