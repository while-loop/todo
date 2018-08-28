package vcs

import (
	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/log"
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

func NewManager(config *config.VcsConfig, router *mux.Router, creator issue.Creator, logger log.AnalysisLogger) *Manager {
	m := &Manager{
		config:   config,
		services: map[string]RepositoryService{},
	}

	m.initServices(creator, logger)
	m.initRouter(router)
	return m
}

func (m *Manager) Services() map[string]RepositoryService {
	return m.services
}

func (m *Manager) initServices(creator issue.Creator, logger log.AnalysisLogger) {
	conf := m.config

	if conf.Github != nil {
		srv := github.NewService(conf.Github, creator, logger)
		m.services[srv.Name()] = srv
	}

	if conf.Gitlab != nil {
		srv := gitlab.NewService(conf.Gitlab, creator, logger)
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
