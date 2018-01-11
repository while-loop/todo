package config

type TrackerConfig struct {
	Github *GithubConfig `json:"github" yaml:"github"`
	Jira   *JiraConfig   `json:"jira" yaml:"jira"`
}

type GithubConfig struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
}

type JiraConfig struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
}
