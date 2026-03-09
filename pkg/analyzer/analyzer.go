package analyzer

import (
	"go/ast"

	"github.com/TSunset/golint/pkg/logger"
	"github.com/TSunset/golint/pkg/rules"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages in slog and zap calls",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
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

			for _, issue := range rules.CheckMessage(logCall.Message) {
				pass.Reportf(logCall.MessagePos, "loglint: %s", issue)
			}

			return true
		})
	}

	return nil, nil
}
