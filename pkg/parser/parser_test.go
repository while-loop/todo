package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDownloadTwoFiles(t *testing.T) {
	a := require.New(t)
	urls := []string{"https://github.com/while-loop/todo/raw/08b3e2fad64e54c061d1ba6324a382e968212a6c/pkg/vcs/github/event_push_test.go",
		"https://github.com/ansible/ansible/raw/781fd7099a0278d3d91557b94da1083f19fad329/test/legacy/roles/test_gce_labels/tasks/test.yml"}

	p := New()
	issues := p.GetTodos(urls...)

	a.Equal(2, len(issues))
	// a.Equal("test when parser has been impl", issues[0].Title) todo
}

func TestDownloadWhenServerIsNotReachable(t *testing.T) {
	a := require.New(t)
	urls := []string{"https://gitakjshdgfasjhgdfhub.com"}

	p := New()
	issues := p.GetTodos(urls...)

	a.Equal(0, len(issues))
}

func TestParsingUnsupportedExtension(t *testing.T) {
	a := require.New(t)

	// readme.md has valid todo comments, but md is an supp extension atm
	urls := []string{"https://github.com/while-loop/todo/raw/cc6b554cccfd3598f6b6342d69c78abcbc5d0128/README.md"}

	p := New()
	issues := p.GetTodos(urls...)

	a.Equal(0, len(issues))
}
