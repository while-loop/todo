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
		iss     tracker.Issue
	}{
		{"// todo hello world", slashRegex, true, tracker.Issue{Title: "hello world"}},
		{"line of code", slashRegex, false, e},
		{"# todo(snake) impl python", hashRegex, true, tracker.Issue{Title: "impl python", Assignee: "snake"}},
		{`// fmt.Println("uh oh") todo(snake) eol comment`, slashRegex, true, tracker.Issue{Title: "eol comment", Assignee: "snake"}},
	}

	for idx, tt := range tcs {
		t.Run(strconv.Itoa(idx), func(inner *testing.T) {
			_, found := parseLine(tt.rexp, tt.comment)

			require.Equal(inner, found, tt.found)
			// require.Equal(inner, iss, tt.iss) todo
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
