package parser

import (
	"regexp"
	"strconv"
	"testing"

	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/while-loop/todo/pkg/issue"
)

func TestRegexInit(t *testing.T) {
	require.Contains(t, slashRegex.String(), `//.*todo`)
	require.Contains(t, hashRegex.String(), `#.*todo`)
}

func TestLangSorted(t *testing.T) {
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
		is      *issue.Issue
	}{
		{"// todo hello world", slashRegex, &issue.Issue{Title: "hello world", Mentions: []string{}, Labels: []string{}, Extras: map[string]interface{}{}}},
		{"line of code", slashRegex, nil},
		{`// using the command "grep -nri -E "^.*//\s*todo.*" > etcdTodos.txt"`, slashRegex, nil},
		{"// this is a list of TODOs found in the root directory of https://github.com/coreos/etcd", slashRegex, nil},
		{"# todo(snake): impl python", hashRegex, &issue.Issue{Title: "impl python", Assignee: "snake", Mentions: []string{}, Labels: []string{}, Extras: map[string]interface{}{}}},
		{"# todo(snake) impl python", hashRegex, &issue.Issue{Title: "impl python", Assignee: "snake", Mentions: []string{}, Labels: []string{}, Extras: map[string]interface{}{}}},
		{`// fmt.Println("uh oh") todo(snake): eol comment`, slashRegex, &issue.Issue{Title: "eol comment", Assignee: "snake", Mentions: []string{}, Labels: []string{}, Extras: map[string]interface{}{}}},
		{"// todo(while-loop): @dev1 finish tests +test +api", slashRegex, &issue.Issue{
			Title:    "finish tests",
			Labels:   []string{"test", "api"},
			Mentions: []string{"@dev1"},
			Assignee: "while-loop",
			Extras:   map[string]interface{}{},
		}},
		{`// todo(while-loop): add ignore keyword to yml config (ParseFile will be a todoParser func)`, slashRegex, &issue.Issue{
			Assignee: "while-loop",
			Title:    "add ignore keyword to yml config (ParseFile will be a todoParser func)",
			Mentions: []string{}, Labels: []string{},
			Extras: map[string]interface{}{},
		}},
		{`// code snippet https://github.com/while-loop/todo/blob/cc6b554cccfd3598f6b6342d69c78abcbc5d0128/pkg/app.go#L17-L25`, slashRegex, nil},
		{`// footer  ###### This issue was generated by [todo](https://github.com/while-loop/todo) on behalf of %s.`, slashRegex, nil},
		{`// foreach file, get all todos`, slashRegex, nil},
		{`	Extras         context.Context // todo change Issue.Extras from Context to map`, slashRegex, nil},
	}

	fmt.Println(slashRegex)
	for idx, tc := range tcs {
		t.Run(strconv.Itoa(idx), func(inner *testing.T) {
			is, _ := parseLine(tc.rexp, tc.comment)
			require.Equal(inner, tc.is, is)
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
		{`// var 3 int todo find all`, []string{}},
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
