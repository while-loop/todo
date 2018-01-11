package vcs

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/tracker"
	"github.com/while-loop/todo/pkg/vcs/config"
	"github.com/while-loop/todo/pkg/vcs/github"
	"github.com/while-loop/todo/pkg/vcs/gitlab"
)

type RepositoryService interface {
	GetRepository(user, project string) error
	Name() string
	Handler() http.Handler
}

type Manager struct {
	config    *config.VcsConfig
	services  map[string]RepositoryService
	issueChan chan tracker.Issue
}

func NewManager(config *config.VcsConfig) *Manager {
	m := &Manager{
		config:    config,
		services:  map[string]RepositoryService{},
		issueChan: make(chan tracker.Issue),
	}

	m.initServices()
	return m
}

func (m *Manager) Services() map[string]RepositoryService {
	return m.services
}

func (m *Manager) IssueChan() <-chan tracker.Issue {
	return m.issueChan
}

func (m *Manager) initServices() {
	conf := m.config

	if conf.Github != nil {
		srv := github.NewService(conf.Github, m.issueChan)
		m.services[srv.Name()] = srv
	}

	if conf.Gitlab != nil {
		srv := gitlab.NewService(conf.Gitlab, m.issueChan)
		m.services[srv.Name()] = srv
	}
}

func (m *Manager) ApplyRouter(router *mux.Router) {
	for _, srvc := range m.services {
		router.NewRoute().Handler(srvc.Handler())
	}
}
