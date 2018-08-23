package vcs

import (
	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/vcs/config"
	"github.com/while-loop/todo/pkg/vcs/github"
	"github.com/while-loop/todo/pkg/vcs/gitlab"
)

type RepositoryService interface {
	Name() string
	Init(webhookRouter *mux.Router)
}

type Manager struct {
	config   *config.VcsConfig
	services map[string]RepositoryService
}

func NewManager(config *config.VcsConfig, router *mux.Router, creator issue.Creator) *Manager {
	m := &Manager{
		config:   config,
		services: map[string]RepositoryService{},
	}

	m.initServices(creator)
	m.initRouter(router)
	return m
}

func (m *Manager) Services() map[string]RepositoryService {
	return m.services
}

func (m *Manager) initServices(creator issue.Creator) {
	conf := m.config

	if conf.Github != nil {
		srv := github.NewService(conf.Github, creator)
		m.services[srv.Name()] = srv
	}

	if conf.Gitlab != nil {
		srv := gitlab.NewService(conf.Gitlab, creator)
		m.services[srv.Name()] = srv
	}
}

func (m *Manager) initRouter(router *mux.Router) {
	hookRoute := router.PathPrefix("/webhook")
	if hookRoute.GetError() != nil {
		panic(hookRoute.GetError())
	}

	hookRouter := hookRoute.Subrouter()
	for _, srvc := range m.services {
		srvc.Init(hookRouter)
	}
}
