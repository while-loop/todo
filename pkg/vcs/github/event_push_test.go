package github

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContentUrl(t *testing.T) {
	exp := "https://github.com/while-loop/test/raw/abcezas123/my/first/file.go"
	sha := "abcezas123"
	path := "my/first/file.go"
	repo := "while-loop/test"

	require.Equal(t, exp, contentUrl(repo, sha, path))
}

func TestFoundTodos(t *testing.T) {
	// TODO test when parser has been impl
}
