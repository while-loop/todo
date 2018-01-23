package tracker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/while-loop/todo/pkg/issue"
)

func TestXor(t *testing.T) {
	s1 := []*issue.Issue{{Title: "my title"}}
	s2 := []*issue.Issue{{Title: "my title"}}
	s3 := []*issue.Issue{{Title: "my asd"}} // found todo !todo
	assert.Equal(t, []*issue.Issue{}, xorIssues(s1, s2))
	assert.Equal(t, []*issue.Issue{{Title: "my asd"}}, xorIssues(s1, s3))
}
