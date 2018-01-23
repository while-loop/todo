package tracker

import (
	"context"

	"time"

	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/log"
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
	issueCh  <-chan []*issue.Issue
}

func NewManager(config *config.TrackerConfig, issueCh <-chan []*issue.Issue) *Manager {
	m := &Manager{
		trackers: map[string]Tracker{},
		config:   config,
		issueCh:  issueCh,
	}

	m.initTrackers()
	go m.loop()
	return m
}

func (m *Manager) loop() {
	for iss := range m.issueCh {
		if len(iss) <= 0 {
			continue
		}

		for _, tr := range m.trackers {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			tIss, err := tr.GetIssues(ctx, iss[0].Owner, iss[0].Repo)
			if err != nil {
				log.Error(err)
				cancel()
				continue
			}

			toCreate := xorIssues(tIss, iss)
			log.Info(iss)
			log.Info(toCreate)
			log.Infof("need to create %d issues", len(toCreate))
			ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
			for _, c := range toCreate {
				log.Info("tocreat: ", c)
				go func() {
					if is, err := tr.CreateIssue(ctx, c); err != nil {
						log.Error(err)
					} else {
						log.Infof("Created issue: %s/%s/%s", is.Owner, is.Repo, is.ID)
					}
				}()
			}
		}
	}
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
