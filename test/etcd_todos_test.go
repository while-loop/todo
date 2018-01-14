package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/while-loop/todo/pkg/parser"
	"github.com/while-loop/todo/pkg/tracker"
	"os"
	"testing"
)

func TestCoreOSEtcd(t *testing.T) {
	// this is a list of TODOs found in the root directory of https://github.com/coreos/etcd
	// using the command `grep -nri -E "^.*//\s*todo.*" > etcdTodos.txt`
	// empty todos are not treated as issues
	fileName := "etcdTodos.java"
	f, err := os.Open(fileName)
	assert.NoError(t, err)

	issues, err := parser.ParseFile(fileName, f)
	assert.NoError(t, err)

	assert.Equal(t, 441, len(issues))
	for _, ei := range expIssues {
		assert.Equal(t, ei, issues[ei.Line-1])
	}

}

// expIssues has a select of todos given in the file
// more will be added as time goes on
var expIssues = []tracker.Issue{
	{Title: "save_fsync latency?", File: "etcdTodos.java", Line: 1},
	{Title: "ensure the entries are continuous and", Assignee: "xiangli", File: "etcdTodos.java", Line: 25},
	{Title: "The original rationale of passing in a pre-allocated", Assignee: "beorn7", File: "etcdTodos.java", Line: 133},
	{Title: "consider a more generally-known optimization for reflect.Value ==> Interface", File: "etcdTodos.java", Line: 200},
	{Title: "support bincUnicodeOther (for now, just use string or bytearray)", File: "etcdTodos.java", Line: 201},
}
