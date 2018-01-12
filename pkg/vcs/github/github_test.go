package github

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
