package config

type VcsConfig struct {
	Github *GithubConfig `json:"github" yaml:"github"`
	Gitlab *GitlabConfig `json:"gitlab" yaml:"gitlab"`
}

type GithubConfig struct {
	AccessToken   string `json:"access_token" yaml:"access_token"`
	PrivateKey    string `json:"private_key" yaml:"private_key"`
	WebhookSecret string `json:"webhook_secret" yaml:"webhook_secret"`
	IssueNumber   int    `json:"issue_number" yaml:"issue_number"`
}

type GitlabConfig struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
}
