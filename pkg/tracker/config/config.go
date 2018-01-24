package config

type TrackerConfig struct {
	Github *GithubConfig `json:"github" yaml:"github"`
	Jira   *JiraConfig   `json:"jira" yaml:"jira"`
}

type GithubConfig struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
	PrivateKey  string `json:"private_key" yaml:"private_key"`
	IssueNumber int    `json:"issue_number" yaml:"issue_number"`
}

type JiraConfig struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
}
