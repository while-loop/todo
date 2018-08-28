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

	if err := s.logger.LogInstallation(event.GetAction(), event.GetInstallation().Account.GetLogin()); err != nil {
		log.Error("err logging installation:", err)
	}
}

func (s *Service) handleRepoInstallation(w http.ResponseWriter, r *http.Request) {
	var event github.InstallationRepositoriesEvent
	if err := json.Unmarshal(r.Context().Value("body").([]byte), &event); err != nil {
		log.Error("failed to decode github InstallationRepositoriesEvent", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("got handleInstallation: %s", render.Render(event))

	var repos []string
	for _, repo := range event.RepositoriesAdded {
		repos = append(repos, repo.GetName())
	}

	for _, repo := range event.RepositoriesRemoved {
		repos = append(repos, repo.GetName())
	}

	if err := s.logger.LogRepoInstallation(event.GetAction(), event.GetInstallation().GetAccount().GetLogin(), repos); err != nil {
		log.Error("err logging installation:", err)
	}
}
