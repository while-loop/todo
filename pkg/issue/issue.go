package issue

import (
	"encoding/json"
)

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
	Extras      map[string]interface{}
}

func New() *Issue {
	return &Issue{Extras: map[string]interface{}{}}
}

func (i Issue) String() string {
	bs, _ := json.Marshal(i)
	return string(bs)
}

func (i Issue) GetString(key string) string {
	if v, exists := i.Extras[key]; exists {
		return v.(string)
	} else {
		return ""
	}
}

func (i Issue) GetInt(key string) int {
	if v, exists := i.Extras[key]; exists {
		return v.(int)
	} else {
		return 0
	}
}
func (i *Issue) MergeMap(m map[string]interface{}) {
	if i.Extras == nil {
		i.Extras = map[string]interface{}{}
	}

	if m == nil {
		return
	}

	for k, v := range m {
		i.Extras[k] = v
	}
}
