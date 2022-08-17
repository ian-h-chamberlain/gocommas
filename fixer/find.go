package fixer

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"strings"
)

func FindMissingCommas(filename string, src []byte) (positions []token.Position, err error) {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, filename, src, parser.AllErrors)
	if err != nil {
		var scannerErrs scanner.ErrorList
		if errors.As(err, &scannerErrs) {
			filtered := filterCommaErrors(scannerErrs)
			if filtered.Len() > 0 {
				return nil, filtered
			}
		} else {
			return nil, err
		}
	}

	positions = []token.Position{}

	ast.Inspect(
		f,
		func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.CompositeLit:
				pos, isMissing := isCompositeLitMissingComma(src, fset, node)
				if isMissing {
					positions = append(positions, pos)
				}
			case *ast.FuncDecl:
				// TODO

			case *ast.CallExpr:
				// TODO
			}

			// for now always descend to children. This might be optimizable later
			// by bailing out for types that cannot have the children we care about
			// (e.g. import statements, struct defs, etc)
			return true
		})

	return positions, nil
}

func filterCommaErrors(errs scanner.ErrorList) scanner.ErrorList {
	var result scanner.ErrorList

	for _, scanErr := range []*scanner.Error(errs) {
		// sorta hacky, but seems like the best option for now
		if !strings.Contains(scanErr.Msg, "missing ',' before newline") {
			result.Add(scanErr.Pos, scanErr.Msg)
		}
	}

	return result

}

func isCompositeLitMissingComma(
	input []byte,
	fset *token.FileSet,
	lit *ast.CompositeLit,
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
