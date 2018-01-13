package parser

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGetGithubFile(t *testing.T) {
	a := require.New(t)
	u := "https://github.com/while-loop/todo/raw/08b3e2fad64e54c061d1ba6324a382e968212a6c/pkg/vcs/github/github_test.go"
	exp := `package github

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTestValidBody(t *testing.T) {
	secret := "mysecret"
	payload := []byte("mypayload")
	computed := "sha1=57852ac1e3fd8e66063cd9d4cb05ea87355bb0b8"

	require.True(t, validBody(payload, secret, computed))
}
`

	rc, err := DownloadFile(&http.Client{Timeout: 5 * time.Second}, u)
	a.NoError(err)
	defer rc.Close()

	content, err := ioutil.ReadAll(rc)
	a.NoError(err)

	a.Equal(exp, string(content))
}

func TestFileNotFound(t *testing.T) {
	a := require.New(t)
	u := "https://github.com/while-loop/todo/raw/ezas123/pkg/vcs/github/github_test.go"

	rc, err := DownloadFile(&http.Client{Timeout: 5 * time.Second}, u)
	a.Contains(err.Error(), "status code 404")
	a.Nil(rc)
}

func TestServerDown(t *testing.T) {
	a := require.New(t)
	u := "https://fakegithuburlTestServerDown.com/who"

	rc, err := DownloadFile(&http.Client{Timeout: 5 * time.Second}, u)
	a.Contains(err.Error(), "connection refused")
	a.Nil(rc)
}
