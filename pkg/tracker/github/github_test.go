package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/while-loop/todo/pkg/issue"
)

func TestBlobUrl(t *testing.T) {
	exp := `https://github.com/while-loop/test/blob/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java#L1-L4`
	assert.Equal(t, exp, createBlobUrl(&issue.Issue{
		Commit: "4912390cee5f96b24e349ae222f3bc25da2708c1",
		Owner:  "while-loop",
		Repo:   "test",
		File:   "https://raw.githubusercontent.com/while-loop/test/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java",
	}))
}
