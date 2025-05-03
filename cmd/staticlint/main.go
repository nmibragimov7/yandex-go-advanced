package main

import (
	"strings"
	"yandex-go-advanced/cmd/staticlint/analyzers/noosexit"

	gocritic "github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

// Обернем OsExitAnalyzer, чтобы он игнорировал файлы из testdata
func wrapOsExitAnalyzer(a *analysis.Analyzer) *analysis.Analyzer {
	wrapped := *a
	oldRun := a.Run
	wrapped.Run = func(pass *analysis.Pass) (interface{}, error) {
		if strings.Contains(pass.Pkg.Path(), "testdata") {
			return nil, nil
		}
		return oldRun(pass)
	}
	return &wrapped
}

// Содержит все анализатора класса SA, анализатор simple и gocritic
// Так же включены все анализаторы passes
// go critic и staticcheck публичные анализаторы
func main() {
	checks := map[string]bool{
		"SA1000": true, "SA1001": true, "SA1002": true, "SA1003": true, "SA1004": true,
		"SA1005": true, "SA1006": true, "SA1007": true, "SA1008": true, "SA1010": true,
		"SA1011": true, "SA1012": true, "SA1013": true, "SA1014": true, "SA1015": true,
		"SA1016": true, "SA1017": true, "SA1018": true, "SA1019": true, "SA1020": true,
		"SA1021": true, "SA1023": true, "SA1024": true, "SA1025": true, "SA1026": true,
		"SA1027": true, "SA1028": true, "SA1029": true, "SA2000": true, "SA2001": true,
		"SA2002": true, "SA2003": true, "SA2004": true, "SA2005": true, "SA2006": true,
		"SA2007": true, "SA2008": true, "SA2009": true, "SA2010": true, "SA2011": true,
		"SA2012": true, "SA2013": true, "SA2014": true, "SA2015": true, "SA2016": true,
		"SA2017": true, "SA2018": true, "SA2019": true, "SA2020": true, "SA2021": true,
		"SA2022": true, "SA2023": true, "SA3000": true, "SA3001": true, "SA4000": true,
		"SA4001": true, "SA4003": true, "SA4004": true, "SA4005": true, "SA4006": true,
		"SA4008": true, "SA4009": true, "SA4010": true, "SA4011": true, "SA4012": true,
		"SA4013": true, "SA4014": true, "SA4015": true, "SA4016": true, "SA4017": true,
		"SA4018": true, "SA4019": true, "SA4020": true, "SA4021": true, "SA4022": true,
		"SA4023": true, "SA5000": true, "SA5001": true, "SA5002": true, "SA5003": true,
		"SA5004": true, "SA5005": true, "SA5007": true, "SA5008": true, "SA5009": true,
		"SA5010": true, "SA5011": true, "SA5012": true, "SA6000": true, "SA6001": true,
		"SA6002": true, "SA6003": true, "SA6005": true, "SA6006": true, "SA9001": true,
		"SA9002": true, "SA9003": true, "SA9004": true, "SA9005": true, "SA9006": true,
		"SA9007": true, "SA9008": true,
	}
	mychecks := []*analysis.Analyzer{
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		shadow.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		waitgroup.Analyzer,
		usesgenerics.Analyzer,
		wrapOsExitAnalyzer(noosexit.OsExitAnalyzer),
		simple.Analyzers[0].Analyzer,
		gocritic.Analyzer,
	}
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}
