package tracker

import (
	"context"

	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/tracker/config"
	"github.com/while-loop/todo/pkg/tracker/github"
)

type Tracker interface {
	GetIssues(ctx context.Context, owner, repo string) ([]*issue.Issue, error)
	CreateIssue(ctx context.Context, issue *issue.Issue) (*issue.Issue, error)
	DeleteIssue(ctx context.Context, issue *issue.Issue) error
	Name() string
}

type Manager struct {
	trackers map[string]Tracker
	config   *config.TrackerConfig
}

func NewManager(config *config.TrackerConfig) *Manager {
	m := &Manager{
		trackers: map[string]Tracker{},
		config:   config,
	}

	m.initTrackers()
	return m
}

func (m *Manager) Trackers() map[string]Tracker {
	return m.trackers
}

func (m *Manager) initTrackers() {
	conf := m.config

	if conf.Github != nil {
		g := github.NewTracker(m.config.Github)
		m.trackers[g.Name()] = g
	}

	if conf.Jira != nil {

	}
}
