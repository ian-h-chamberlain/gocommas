package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	src, err := os.ReadFile("testdata/missing_trailing.go")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse src but stop after processing the imports.
	f, _ := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	added := []token.Position{}

	ast.Inspect(
		f,
		func(node ast.Node) bool {
			switch node.(type) {
			case *ast.BadDecl, *ast.BadExpr, *ast.BadStmt:
				fmt.Println(node.Pos(), node.End())
			case *ast.CompositeLit:
				pos, isMissing := isMissingTrailingComma(src, fset, node.(*ast.CompositeLit))
				if isMissing {
					added = append(added, pos)
				}
			}
			return true
		})

	offset := 0
	for _, commaPos := range added {
		fmt.Println("missing trailing commma:", commaPos)

		// insert a comma at the given position
		i := commaPos.Offset + offset

		fmt.Println("inserting comma at", i)

		after := append([]byte(","), src[i:]...)
		src = append(src[:i], after...)
		offset += 1
	}

	fmt.Println("-----------------------------------")
	fmt.Println(string(src))
}

func isMissingTrailingComma(
	input []byte, fset *token.FileSet, lit *ast.CompositeLit,
) (token.Position, bool) {
	if len(lit.Elts) == 0 {
		return token.Position{}, false
	}

	lastEl := lit.Elts[len(lit.Elts)-1]

	lastElPos := fset.Position(lastEl.End())
	rbracePos := fset.Position(lit.Rbrace)

	if rbracePos.Line <= lastElPos.Line {
		// trailing comma not required if brace is on same line
		return token.Position{}, false
	}

	nextLine := fset.File(lastEl.Pos()).LineStart(lastElPos.Line + 1)
	if !nextLine.IsValid() {
		return token.Position{}, false
	}
	nextLinePos := fset.Position(nextLine)

	searchString := input[lastElPos.Offset:nextLinePos.Offset]
	if !bytes.ContainsRune(searchString, ',') {
		return lastElPos, true
	}

	return token.Position{}, false
}
