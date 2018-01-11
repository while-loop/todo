package tracker

import (
	"encoding/json"

	"github.com/while-loop/todo/pkg/tracker/config"
)

type Tracker interface {
	GetIssues() ([]*Issue, error)
	CreateIssue(issue *Issue) (*Issue, error)
	DeleteIssue(issue *Issue) error
}

type Issue struct {
	ID          string
	Title       string
	Description string
	Assignee    string
	Author      string
	Mentions    string
	Labels      []string
}

func (i *Issue) String() string {
	bs, _ := json.Marshal(i)
	return string(bs)
}

type Manager struct {
	trackers map[string]Tracker
	config   *config.TrackerConfig
}

func NewManager(config *config.TrackerConfig) *Manager {
	m := &Manager{
		trackers: map[string]Tracker{},
	}

	return m
}

func (m *Manager) Trackers() map[string]Tracker {
	return m.trackers
}

func (m *Manager) initPublishers() {
	conf := m.config

	if conf.Github != nil {

	}

	if conf.Jira != nil {

	}
}
