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

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil, false
	}

	msg, err := strconv.Unquote(lit.Value)
	if err != nil {
		return nil, false
	}

	loggerKind, ok := detectLoggerKind(pass, sel)
	if !ok {
		return nil, false
	}

	return &LogCall{
		Call:       call,
		Message:    msg,
		MessagePos: lit.Pos(),
		MessageEnd: lit.End(),
		Method:     sel.Sel.Name,
		LoggerKind: loggerKind,
	}, true
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
