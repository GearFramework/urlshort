package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMyAnalyzerNoMain(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitAnalyzer, "./testpkgnomain")
}