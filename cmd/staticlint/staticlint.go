// Package multichecker for check call os.Exist() in main()
package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

// OsExitAnalyzer instance of analyzer
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for call os.Exit() from main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	isModuleOS := func(c *ast.SelectorExpr) bool {
		if idn, ok := c.X.(*ast.Ident); ok {
			fmt.Println("Module: ", idn.Name)
			return idn.Name == "os"
		}
		return false
	}
	callFunc := func(c *ast.CallExpr) {
		if s, ok := c.Fun.(*ast.SelectorExpr); ok {
			if isModuleOS(s) && s.Sel.Name == "Exit" {
				fmt.Println("Func: ", s.Sel.Name)
				pass.Reportf(c.Pos(), "directly call os.Exit")
			}
		}
	}
	declFuncMain := func(c *ast.FuncDecl) {
		if c.Name.Name != "main" {
			return
		}
		ast.Inspect(c.Body, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				callFunc(x)
			}
			return true
		})
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if file.Name.Name == "main" {
				switch x := node.(type) {
				case *ast.FuncDecl:
					declFuncMain(x)
				}
			}
			return true
		})
	}
	return nil, nil
}

func main() {
	multichecker.Main(
		OsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	)
}
