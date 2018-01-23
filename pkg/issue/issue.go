package issue

import "encoding/json"

type Issue struct {
	ID          string
	Title       string
	Description string
	Assignee    string
	Author      string
	Mentions    []string
	Labels      []string
	File        string
	Line        int
	Owner       string
	Repo        string
	Commit      string
}

func (i *Issue) String() string {
	bs, _ := json.Marshal(i)
	return string(bs)
}
