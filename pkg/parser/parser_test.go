package parser

import (
	"testing"

	"reflect"

	"context"
	"net/http"

	"github.com/stretchr/testify/require"
	"github.com/while-loop/todo/pkg/issue"
)

// uncomment test
func testDownloadTwoFiles(t *testing.T) {
	a := require.New(t)
	urls := []string{"https://github.com/while-loop/todo/raw/08b3e2fad64e54c061d1ba6324a382e968212a6c/pkg/vcs/github/event_push_test.go",
		"https://github.com/ansible/ansible/raw/781fd7099a0278d3d91557b94da1083f19fad329/test/legacy/roles/test_gce_labels/tasks/test.yml"}

	iss := []*issue.Issue{
		{Assignee: "erjohnso", Title: "write more tests", File: urls[1], Line: 28, Mentions: []string{}, Labels: []string{}, Commit: "master"},
		{Title: "test when parser has been impl", File: urls[0], Line: 18, Mentions: []string{}, Labels: []string{}, Commit: "master"},
	}

	p := New()
	issues := p.GetTodos(context.Background(), http.DefaultClient, urls...)
	a.Equal(len(iss), len(issues))
	ttl := 0 // keep a count of total issues matched. since order of received issues are random due to goroutines

	for _, expIs := range iss {
		for i, actIs := range issues {
			if reflect.DeepEqual(expIs, actIs) {
				ttl++
				issues = remove(issues, i)
				break
			}
		}
	}

	a.Equal(len(iss), ttl)
}

func TestDownloadWhenServerIsNotReachable(t *testing.T) {
	a := require.New(t)
	urls := []string{"https://gitakjshdgfasjhgdfhub.com"}

	p := New()
	issues := p.GetTodos(context.Background(), http.DefaultClient, urls...)

	a.Equal(0, len(issues))
}

func TestParsingUnsupportedExtension(t *testing.T) {
	a := require.New(t)

	// readme.md has valid todo comments, but md is an supp extension atm
	urls := []string{"https://github.com/while-loop/todo/raw/cc6b554cccfd3598f6b6342d69c78abcbc5d0128/README.md"}

	p := New()
	issues := p.GetTodos(context.Background(), http.DefaultClient, urls...)

	a.Equal(0, len(issues))
}

func remove(s []*issue.Issue, i int) []*issue.Issue {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
