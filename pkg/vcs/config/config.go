package config

type VcsConfig struct {
	Github *GithubConfig `json:"github" yaml:"github"`
	Gitlab *GitlabConfig `json:"gitlab" yaml:"gitlab"`
}

type GithubConfig struct {
	AccessToken   string `json:"access_token" yaml:"access_token"`
	WebhookSecret string `json:"webhook_secret" yaml:"webhook_secret"`
}

type GitlabConfig struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
}
