package github

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContentUrl(t *testing.T) {
	exp := "https://raw.githubusercontent.com/while-loop/test/abcezas123/my/first/file.go"
	sha := "abcezas123"
	path := "my/first/file.go"
	owner := "while-loop"
	repo := "test"

	require.Equal(t, exp, contentUrl(owner, repo, sha, path))
}

func TestFoundTodos(t *testing.T) {
	// TODO test when parser has been impl
}
