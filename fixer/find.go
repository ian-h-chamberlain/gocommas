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

	fileRoot, err := parser.ParseFile(fset, filename, src, parser.AllErrors)
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

	finder := missingCommaFinder{Src: src, Fset: fset}

	ast.Inspect(fileRoot, finder.VisitNode)

	return finder.Positions, nil
}

func filterCommaErrors(errs scanner.ErrorList) scanner.ErrorList {
	var result scanner.ErrorList

	// We may want to keep track of the comma errors for later, so we can exit
	// with error if there was something we couldn't fix. Basically just ensure
	// all the positions are accounted for.

	// Hmm, actually I guess we could just use the errors directly to find
	// positions, but it might be more robust to use the AST and just check our
	// work using the errors. Let's keep trying it that way for now

	for _, scanErr := range []*scanner.Error(errs) {
		// sorta hacky, but seems like the best option for now
		if !strings.Contains(scanErr.Msg, "missing ',' before newline") {
			result.Add(scanErr.Pos, scanErr.Msg)
		}
	}

	return result

}

type missingCommaFinder struct {
	Src       []byte
	Fset      *token.FileSet
	Positions []token.Position
}

func (f *missingCommaFinder) VisitNode(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.FieldList:
		f.Positions = append(f.Positions, f.findInFieldList(node)...)
	case *ast.CompositeLit:
		f.Positions = append(f.Positions, f.findInCompositeLit(node)...)
	case *ast.FuncDecl:
		// all child types of FuncDecl should be covered by FieldList
	case *ast.CallExpr:
		// TODO
	}

	// for now always descend to children. This might be optimizable later
	// by bailing out for types that cannot have the children we care about
	// (e.g. import statements, struct defs, etc)
	return true

}

func (f *missingCommaFinder) findInCompositeLit(lit *ast.CompositeLit) []token.Position {
	if pos, ok := findMissingComma(
		f.Src,
		f.Fset,
		lit.Elts,
		lit.Rbrace,
	); ok {
		return []token.Position{pos}
	}

	return nil
}

func (f *missingCommaFinder) findInFuncDecl(funcDecl *ast.FuncDecl) []token.Position {
	positions := []token.Position{}

	if funcDecl.Recv != nil {
		if pos, ok := findMissingComma(
			f.Src,
			f.Fset,
			funcDecl.Recv.List,
			funcDecl.Recv.Closing,
		); ok {
			positions = append(positions, pos)
		}
	}

	if funcDecl.Type != nil {
		if pos, ok := findMissingComma(
			f.Src,
			f.Fset,
			funcDecl.Type.Params.List,
			funcDecl.Type.Params.Closing,
		); ok {
			positions = append(positions, pos)
		}

		if funcDecl.Type.Results != nil {
			if pos, ok := findMissingComma(
				f.Src,
				f.Fset,
				funcDecl.Type.Results.List,
				funcDecl.Type.Results.Closing,
			); ok {
				positions = append(positions, pos)
			}
		}
	}

	return positions
}

func (f *missingCommaFinder) findInFieldList(fieldList *ast.FieldList) []token.Position {
	if pos, ok := findMissingComma(
		f.Src,
		f.Fset,
		fieldList.List,
		fieldList.Closing,
	); ok {
		return []token.Position{pos}
	}

	return nil
}

func findMissingComma[Elem ast.Node](
	input []byte,
	fset *token.FileSet,
	nodes []Elem,
	closing token.Pos,
) (token.Position, bool) {
	if len(nodes) == 0 {
		return token.Position{}, false
	}

	lastEl := nodes[len(nodes)-1]

	lastElPos := fset.Position(lastEl.End())
	rbracePos := fset.Position(closing)

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
