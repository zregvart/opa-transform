package main

import (
	"context"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

const rule = "rule.rego"

func main() {
	src, err := os.ReadFile(rule)
	if err != nil {
		panic(err)
	}

	mod := ast.MustParseModule(string(src))

	var body ast.Body
	var exprIdx = 0
	ast.NewGenericVisitor(func(x any) bool {
		if b, ok := x.(ast.Body); ok {
			body = b
			exprIdx = 0
		}
		if _, ok := x.(*ast.Expr); ok {
			exprIdx++
		}
		if c, ok := x.(ast.Call); ok && c[0].Value.String() == "custom_function" {
			body[exprIdx-1] = ast.MustParseExpr(`val = ""`) // TODO smarter replacement
		}
		return false
	}).Walk(mod)

	fmt.Println("Modified AST")
	fmt.Println("---------")
	fmt.Println(mod)
	fmt.Println("---------")

	compiler := ast.NewCompiler()
	compiler.Compile(map[string]*ast.Module{
		rule: mod,
	})

	if len(compiler.Errors) != 0 {
		for _, err := range compiler.Errors {
			fmt.Println(err)
		}

		return
	}

	result, err := rego.New(rego.Compiler(compiler), rego.Query("data.rule.custom")).Eval(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Result")
	fmt.Println("---------")
	fmt.Println(result)
	fmt.Println("---------")
}
