package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"strings"

	"github.com/derision-test/go-mockgen/internal"
)

// archive holds information about a library and its corresponding archive.
//
// It is mainly used as a flag through the cff CLI. Its format closely follows
// the format used in Bazel rules_go: https://github.com/bazelbuild/rules_go/blob/8ea79bbd5e6ea09dc611c245d1dc09ef7ab7118a/go/private/actions/compile.bzl#L20
//
// The following is the flag format:
//
//	--archive=IMPORTPATHS=IMPORTMAP=FILE=EXPORT
//
// For example,
//
//	--archive=github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go
//
// However, we only use the ImportMap and File attribute from this flag. In the
// future, we may use ImportPaths to resolve import aliases.
type archive struct {
	// ImportMap refers to the actual import path to the library this archive
	// represents. While the naming may be confusing, this closely follows Bazel
	// rules_go conventions.
	//
	// See https://github.com/bazelbuild/rules_go/blob/f7a8cb6b9158006e5dfc91074f9636820a446921/go/core.rst#go_library.
	ImportMap string
	File      string
}

// parseArchive parses the archive string to the internal.Archive type.
//
// The following is the flag format:
//
//	--archive=IMPORTPATHS=IMPORTMAP=FILE=EXPORT
//
// For example,
//
//	--archive=github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go
//
// The flag is structured in this format to closely follow https://github.com/bazelbuild/rules_go/blob/8ea79bbd5e6ea09dc611c245d1dc09ef7ab7118a/go/private/actions/compile.bzl#L20;
// however, the IMPORTPATHS and EXPORT elements are ignored. There may be future
// work involved in resolving import aliases, using IMPORTPATHS.
func parseArchive(a string) (archive, error) {
	args := strings.Split(a, "=")
	if len(args) != 4 {
		return archive{}, fmt.Errorf("expected 4 elements, got %d", len(args))
	}

	// Currently, we ignore the IMPORTPATHS and EXPORT elements.
	return archive{
		ImportMap: args[1],
		File:      args[2],
	}, nil
}

func PackagesArchive(p loadParams) ([]*internal.GoPackage, error) {
	files := make([]*ast.File, 0, len(p.sources))
	for _, src := range p.sources {
		f, err := parser.ParseFile(p.fset, src, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q: %v", src, err)
		}
		files = append(files, f)
	}

	// Build an importer using the imports map built by reading dependency
	// archives, and use it to build the *types.Package and *types.Info for the
	// source files.
	imp, err := newImporter(p.fset, p.archives, p.stdlibRoot)
	if err != nil {
		return nil, err
	}
	conf := types.Config{Importer: imp}
	typesInfo := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	pkg, err := conf.Check(p.importPaths[0], p.fset, files, typesInfo)
	if err != nil {
		return nil, fmt.Errorf("error building pkg %q: %v", p.importPaths[0], err)
	}
	return []*internal.GoPackage{
		{
			PkgPath:         pkg.Path(),
			CompiledGoFiles: p.sources,
			Syntax:          files,
			Types:           pkg,
			TypesInfo:       typesInfo,
		},
	}, nil
}
