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

	// TODO: we might be able to prtially support formatting "fragments" by
	// using parser.ParseExprFrom(), possibly with some heuristic modifications
	// like wrapping it in `func() {}` to support statements, and eventually falling back
	// if nothing works. Alternatively, maybe a "selection mode" CLI option would
	// work a little better than just trying everything every time.

	fileRoot, err := parser.ParseFile(fset, filename, src, parser.AllErrors|parser.SkipObjectResolution)
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
	case *ast.StructType, *ast.InterfaceType:
		// struct definitions contain a FieldList, but should be separated by
		// newlines/semicolons, not commas, so don't descend into these
		return false

	case *ast.FieldList:
		f.Positions = append(f.Positions, f.findInFieldList(node)...)

	case *ast.CompositeLit:
		f.Positions = append(f.Positions, f.findInCompositeLit(node)...)

	case *ast.CallExpr:
		f.Positions = append(f.Positions, f.findInCallExpr(node)...)

	case *ast.FuncDecl:
		// all child types of FuncDecl should be covered by FieldList
	}

	// For now always descend to children. This might be optimizable later
	// by bailing out for types that cannot have the children we care about
	// (e.g. import blocks, etc)
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

func (f *missingCommaFinder) findInCallExpr(callExpr *ast.CallExpr) []token.Position {
	if callExpr.Ellipsis != token.NoPos {
		// gofmt handles trailing ellipsis without any trouble, since a trailing
		// comma is not required after an ellipsis. So if we see an ellipsis in
		// the call expression, we don't try to fix anything.
		return nil
	}

	args := []ast.Node{}
	for _, arg := range callExpr.Args {
		args = append(args, arg)
	}

	if pos, ok := findMissingComma(
		f.Src,
		f.Fset,
		args,
		callExpr.Rparen,
	); ok {
		return []token.Position{pos}
	}

	return nil
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
