package tracker

import (
	"context"

	"sync"

	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/tracker/config"
	"github.com/while-loop/todo/pkg/tracker/github"
)

type Tracker interface {
	GetIssues(ctx context.Context, ref *issue.Issue) ([]*issue.Issue, error)
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

func (m *Manager) Create(issues []*issue.Issue) error {
	if len(issues) < 0 {
		return nil
	}

	for _, tr := range m.trackers {
		tIss, err := tr.GetIssues(context.Background(), issues[0])
		if err != nil {
			log.Error(err)
			continue
		}

		toCreate := xorIssues(tIss, issues)
		log.Info(issues)
		log.Info(toCreate)
		log.Infof("need to create %d issues", len(toCreate))
		var wg sync.WaitGroup
		for _, cr := range toCreate {
			wg.Add(1)
			go func(i *issue.Issue) {
				log.Info("tocreat: ", i)
				if is, err := tr.CreateIssue(context.Background(), i); err != nil {
					log.Error(err)
				} else {
					log.Infof("Created issue: %s/%s/%s", is.Owner, is.Repo, is.ID)
				}
				wg.Done()
			}(cr)

			wg.Wait()
		}
	}

	return nil
}

func xorIssues(repoIssues []*issue.Issue, pushIssues []*issue.Issue) []*issue.Issue {
	toCreate := make([]*issue.Issue, 0)
	for _, pIs := range pushIssues {
		if !contains(pIs, repoIssues) {
			toCreate = append(toCreate, pIs)
		}
	}

	return toCreate
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

func contains(is *issue.Issue, iss []*issue.Issue) bool {
	for _, i := range iss {
		if i.Title == is.Title {
			return true
		}
	}
	return false
}
