package test

import (
	"os"
	"reflect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/parser"
	"testing"
)

var fileName = `etcdTodos.test`

func TestEtd(t *testing.T) {
	// this is a list of TODOs found in the root directory of https://github.com/coreos/etcd
	// using the command `grep -nri -E "^.*//\s*todo.*" > etcdTodos.txt`
	// empty todos are not treated as issues

	f, err := os.Open(fileName)
	require.Nil(t, err)

	issues, lines, err := parser.ParseFile(fileName, f)
	require.Nil(t, err)

	assert.Equal(t, len(issues), lines)
	assert.Equal(t, 441, len(issues))

	for _, ei := range expIssues {
		assert.True(t, reflect.DeepEqual(ei, issues[ei.Line-1]), "title: %s", ei.Title)
	}
}

// expIssues has a select of todos given in the file
// more will be added as time goes on
var expIssues = []*issue.Issue{
	i(issue.Issue{Title: "save_fsync latency?", File: fileName, Line: 1}),
	i(issue.Issue{Title: "ensure the entries are continuous and", Assignee: "xiangli", File: fileName, Line: 25}),
	i(issue.Issue{Title: "The original rationale of passing in a pre-allocated", Assignee: "beorn7", File: fileName, Line: 133}),
	i(issue.Issue{Title: "consider a more generally-known optimization for reflect.Value ==> Interface", File: fileName, Line: 200}),
	i(issue.Issue{Title: "support bincUnicodeOther (for now, just use string or bytearray)", File: fileName, Line: 201}),
}

func i(is issue.Issue) *issue.Issue {
	is.Mentions = []string{}
	is.Labels = []string{}
	is.Commit = "master"
	is.Extras = map[string]interface{}{}
	return &is
}
