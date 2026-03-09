package logger

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

var allowedMethods = map[string]struct{}{
	"Debug": {},
	"Info":  {},
	"Warn":  {},
	"Error": {},
	"Fatal": {},
	"Panic": {},
}

type LogCall struct {
	Call       *ast.CallExpr
	Message    string
	MessagePos token.Pos
	MessageEnd token.Pos
	LoggerKind string
	Method     string
}

func MatchLogCall(pass *analysis.Pass, call *ast.CallExpr) (*LogCall, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	if _, ok := allowedMethods[sel.Sel.Name]; !ok {
		return nil, false
	}

	if len(call.Args) == 0 {
		return nil, false
	}

	msg, ok := extractMessage(call.Args[0])
	if !ok {
		return nil, false
	}

	loggerKind, ok := detectLoggerKind(pass, sel)
	if !ok {
		return nil, false
	}

	return &LogCall{
		Call:       call,
		Message:    msg,
		MessagePos: call.Args[0].Pos(),
		MessageEnd: call.Args[0].End(),
		Method:     sel.Sel.Name,
		LoggerKind: loggerKind,
	}, true
}

func extractMessage(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", false
		}

		msg, err := strconv.Unquote(e.Value)
		if err != nil {
			return "", false
		}

		return msg, true
	case *ast.ParenExpr:
		return extractMessage(e.X)
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", false
		}

		left, leftOK := extractMessage(e.X)
		if !leftOK {
			// If we cannot recover the beginning of the message,
			// skip the call to avoid unreliable diagnostics.
			return "", false
		}

		right, rightOK := extractMessage(e.Y)
		if !rightOK {
			return left, true
		}

		return left + right, true
	default:
		return "", false
	}
}

func detectLoggerKind(pass *analysis.Pass, sel *ast.SelectorExpr) (string, bool) {
	if isSlogPackageCall(pass, sel) {
		return "slog", true
	}

	if kind, ok := detectBySelection(pass, sel); ok {
		return kind, true
	}

	if kind, ok := detectByReceiverType(pass, sel); ok {
		return kind, true
	}

	return "", false
}

func isSlogPackageCall(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.Uses[ident]
	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	return pkgName.Imported().Path() == "log/slog"
}

func detectBySelection(pass *analysis.Pass, sel *ast.SelectorExpr) (string, bool) {
	selection := pass.TypesInfo.Selections[sel]
	if selection == nil {
		return "", false
	}

	obj := selection.Obj()
	if obj == nil || obj.Pkg() == nil {
		return "", false
	}

	switch obj.Pkg().Path() {
	case "log/slog":
		return "slog", true
	case "go.uber.org/zap":
		return "zap", true
	default:
		return "", false
	}
}

func detectByReceiverType(pass *analysis.Pass, sel *ast.SelectorExpr) (string, bool) {
	t := receiverTypeOfSelector(pass, sel)
	if t == nil {
		return "", false
	}

	if isNamedType(t, "log/slog", "Logger") {
		return "slog", true
	}

	if isNamedType(t, "go.uber.org/zap", "Logger") {
		return "zap", true
	}

	return "", false
}

func receiverTypeOfSelector(pass *analysis.Pass, sel *ast.SelectorExpr) types.Type {
	if selection := pass.TypesInfo.Selections[sel]; selection != nil {
		return selection.Recv()
	}

	if tv, ok := pass.TypesInfo.Types[sel.X]; ok && tv.Type != nil {
		return tv.Type
	}

	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj := pass.TypesInfo.Uses[ident]; obj != nil {
			return obj.Type()
		}
	}

	return nil
}

func isNamedType(t types.Type, pkgPath, typeName string) bool {
	if t == nil {
		return false
	}

	for {
		ptr, ok := t.(*types.Pointer)
		if !ok {
			break
		}
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == pkgPath && obj.Name() == typeName
}
