package tracker

type Tracker interface {
	GetIssues() ([]*Issue, error)
	CreateIssue(issue *Issue) (*Issue, error)
	DeleteIssue(issue *Issue) error
}

type Issue struct {
	title string
	description string
	assignee string
	author string
	mentions string
}
