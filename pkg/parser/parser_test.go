package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDownloadTwoFiles(t *testing.T) {
	a := require.New(t)
	urls := []string{"https://github.com/while-loop/todo/raw/08b3e2fad64e54c061d1ba6324a382e968212a6c/pkg/vcs/github/event_push_test.go"}

	p := New()
	issues := p.GetTodos(urls...)

	a.Equal(1, len(issues))
	// a.Equal("test when parser has been impl", issues[0].Title) todo
}
