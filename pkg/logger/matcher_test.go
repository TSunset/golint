package logger_test

import (
	"go/ast"
	"testing"

	"github.com/TSunset/golint/pkg/logger"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

var testAnalyzer = &analysis.Analyzer{
	Name: "loggermatchertest",
	Doc:  "tests logger matcher",
	Run: func(pass *analysis.Pass) (any, error) {
		for _, file := range pass.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				logCall, ok := logger.MatchLogCall(pass, call)
				if !ok {
					return true
				}

				pass.Reportf(logCall.MessagePos, "matched %s", logCall.LoggerKind)
				return true
			})
		}

		return nil, nil
	},
}

func TestMatchLogCall(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, testAnalyzer, "a")
}
