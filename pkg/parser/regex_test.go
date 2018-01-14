package parser

import (
	"github.com/stretchr/testify/require"
	"github.com/while-loop/todo/pkg/tracker"
	"regexp"
	"strconv"
	"testing"
)

func TestRegexInit(t *testing.T) {
	require.Contains(t, slashRegex.String(), `//.*todo`)
	require.Contains(t, hashRegex.String(), `#.*todo`)
}

func TestLangSorted(t *testing.T) {
	// bash < py
	// c < go
	tcs := []struct {
		l1, l2   string
		kind     []string
		lessThan bool
	}{
		{"c", "go", slashLangs, true},
		{"java", "go", slashLangs, false},
		{"bash", "py", hashLangs, true},
		{"yml", "sh", hashLangs, false},
	}

	for i, tc := range tcs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			i1 := getIdx(tc.kind, tc.l1)
			i2 := getIdx(tc.kind, tc.l2)
			if tc.lessThan {
				require.True(t, i1 < i2)
			} else {
				require.True(t, i1 > i2)
			}
		})
	}
}

func TestTODOs(t *testing.T) {
	tcs := []struct {
		comment string
		rexp    *regexp.Regexp
		found   bool
		is      tracker.Issue
	}{
		{"// todo hello world", slashRegex, true, tracker.Issue{Title: "hello world"}},
		{"line of code", slashRegex, false, e},
		{"# todo(snake): impl python", hashRegex, true, tracker.Issue{Title: "impl python", Assignee: "snake"}},
		{`// fmt.Println("uh oh") todo(snake): eol comment`, slashRegex, true, tracker.Issue{Title: "eol comment", Assignee: "snake"}},
		{"// todo(while-loop): @dev1 finish tests +test +api", slashRegex, true, tracker.Issue{
			Title:    "finish tests",
			Labels:   []string{"test", "api"},
			Mentions: []string{"@dev1"},
			Assignee: "while-loop",
		}},
	}

	for idx, tt := range tcs {
		t.Run(strconv.Itoa(idx), func(inner *testing.T) {
			is, found := parseLine(tt.rexp, tt.comment)

			require.Equal(inner, found, tt.found)
			require.Equal(inner, is, tt.is)
		})
	}
}

func TestMentions(t *testing.T) {
	tcs := []struct {
		cmt      string
		mentions []string
	}{
		{`// todo(while-loop): find all @user1 users in db  @wh3n7hi5-ends-i_win`, []string{"@user1", "@wh3n7hi5-ends-i_win"}},
		{`// var 3 int todo find all`, nil},
	}

	for idx, tt := range tcs {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			require.Equal(t, tt.mentions, mentionsRegex.FindAllString(tt.cmt, -1))
		})
	}
}

func TestLabels(t *testing.T) {
	tcs := []struct {
		cmt    string
		labels []string
	}{
		{`// todo(while-loop): find all +users users in db +api/users +wh3n7hi5-ends-i_win`, []string{"users", "api/users", "wh3n7hi5-ends-i_win"}},
		{`// var 3 int todo find all`, nil},
	}

	for idx, tt := range tcs {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			require.Equal(t, tt.labels, parseLabels(tt.cmt))
		})
	}
}

func getIdx(arr []string, item string) int {

	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			return i
		}
	}
	return -1
}

var e = tracker.Issue{}
