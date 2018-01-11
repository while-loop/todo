package publisher

import (
	"github.com/while-loop/todo/pkg/publisher/config"
	"github.com/while-loop/todo/pkg/publisher/mock"
	"github.com/while-loop/todo/pkg/tracker"
)

type Publisher interface {
	Publish(issue *tracker.Issue) error
}

type Manager struct {
	config     *config.PublisherConfig
	publishers map[string]Publisher
}

func NewManager(config *config.PublisherConfig) *Manager {
	m := &Manager{
		config:     config,
		publishers: map[string]Publisher{},
	}

	m.initPublishers()
	return m
}

func (m *Manager) Publishers() map[string]Publisher {
	return m.publishers
}

func (m *Manager) initPublishers() {
	conf := m.config

	if conf.Mock != nil {
		srv := mock.NewPublisher(conf.Mock)
		m.publishers[srv.Name()] = srv
	}
}
