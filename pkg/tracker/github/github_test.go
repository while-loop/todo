package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/while-loop/todo/pkg/issue"
)

func TestBlobUrl(t *testing.T) {
	cases := []struct {
		name  string
		exp   string
		issue issue.Issue
	}{
		{"no line no given", `https://github.com/while-loop/test/blob/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java`, issue.Issue{
			Commit: "4912390cee5f96b24e349ae222f3bc25da2708c1",
			Owner:  "while-loop",
			Repo:   "test",
			File:   "https://raw.githubusercontent.com/while-loop/test/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java",
		}},
		{
			"line, but no total", `https://github.com/while-loop/test/blob/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java#L5`, issue.Issue{
				Commit: "4912390cee5f96b24e349ae222f3bc25da2708c1",
				Owner:  "while-loop",
				Repo:   "test",
				Line:   5,
				File:   "https://raw.githubusercontent.com/while-loop/test/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java",
			}},
		{
			"line and total, top file", `https://github.com/while-loop/test/blob/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java#L1-L7`, issue.Issue{
				Commit: "4912390cee5f96b24e349ae222f3bc25da2708c1",
				Owner:  "while-loop",
				Repo:   "test",
				Line:   1,
				Extras: map[string]interface{}{"lines": 100},
				File:   "https://raw.githubusercontent.com/while-loop/test/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java",
			}},
		{
			"line and total, top file short file", `https://github.com/while-loop/test/blob/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java#L1-L3`, issue.Issue{
				Commit: "4912390cee5f96b24e349ae222f3bc25da2708c1",
				Owner:  "while-loop",
				Repo:   "test",
				Line:   1,
				Extras: map[string]interface{}{"lines": 3},
				File:   "https://raw.githubusercontent.com/while-loop/test/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java",
			}},
		{
			"line and total, middle file", `https://github.com/while-loop/test/blob/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java#L47-L53`, issue.Issue{
				Commit: "4912390cee5f96b24e349ae222f3bc25da2708c1",
				Owner:  "while-loop",
				Repo:   "test",
				Line:   50,
				Extras: map[string]interface{}{"lines": 100},
				File:   "https://raw.githubusercontent.com/while-loop/test/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java",
			}},
		{
			"line and total, end file", `https://github.com/while-loop/test/blob/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java#L48-L53`, issue.Issue{
				Commit: "4912390cee5f96b24e349ae222f3bc25da2708c1",
				Owner:  "while-loop",
				Repo:   "test",
				Line:   51,
				Extras: map[string]interface{}{"lines": 53},
				File:   "https://raw.githubusercontent.com/while-loop/test/4912390cee5f96b24e349ae222f3bc25da2708c1/test.java",
			}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exp, createBlobUrl(&tc.issue))
		})
	}
}
