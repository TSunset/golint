package plugin

import (
	"github.com/TSunset/golint/pkg/analyzer"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

type Plugin struct{}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func New(conf any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

func init() {
	register.Plugin("loglint", New)
}
