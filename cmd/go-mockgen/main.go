package main

import (
	"errors"
	"fmt"
	"go/token"
	"log"
	"strings"

	"github.com/derision-test/go-mockgen/internal"
	"github.com/derision-test/go-mockgen/internal/mockgen/generation"
	"github.com/derision-test/go-mockgen/internal/mockgen/types"
	"golang.org/x/tools/go/packages"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("go-mockgen: ")
}

func main() {
	if err := mainErr(); err != nil {
		message := fmt.Sprintf("error: %s\n", err.Error())

		if solvableError, ok := err.(solvableError); ok {
			message += "\nPossible solutions:\n"

			for _, hint := range solvableError.Solutions() {
				message += fmt.Sprintf("  - %s\n", hint)
			}

			message += "\n"
		}

		log.Fatalf(message)
	}
}

type solvableError interface {
	Solutions() []string
}

func mainErr() error {
	allOptions, err := parseAndValidateOptions()
	if err != nil {
		return err
	}

	var importPaths []string
	for _, opts := range allOptions {
		for _, packageOpts := range opts.PackageOptions {
			importPaths = append(importPaths, packageOpts.ImportPaths...)
		}
	}

	archives := make([]archive, 0, len(allOptions[0].PackageOptions[0].Archives))
	for _, archive := range allOptions[0].PackageOptions[0].Archives {
		a, err := parseArchive(archive)
		if err != nil {
			return err
		}
		archives = append(archives, a)
	}

	log.Printf("loading data for %d packages\n", len(importPaths))

	pkgs, err := loadPackages(loadParams{
		fset:        token.NewFileSet(),
		importPaths: importPaths,
		// gcexportdata
		archives:   archives,
		sources:    allOptions[0].PackageOptions[0].Sources,
		stdlibRoot: allOptions[0].PackageOptions[0].StdlibRoot,
	})
	if err != nil {
		return fmt.Errorf("could not load packages %s (%s)", strings.Join(importPaths, ","), err.Error())
	}

	for _, opts := range allOptions {
		typePackageOpts := make([]types.PackageOptions, 0, len(opts.PackageOptions))
		for _, packageOpts := range opts.PackageOptions {
			typePackageOpts = append(typePackageOpts, types.PackageOptions(packageOpts))
		}

		ifaces, err := types.Extract(pkgs, typePackageOpts)
		if err != nil {
			return err
		}

		nameMap := make(map[string]struct{}, len(ifaces))
		for _, t := range ifaces {
			nameMap[strings.ToLower(t.Name)] = struct{}{}
		}

		for _, packageOpts := range opts.PackageOptions {
			for _, name := range packageOpts.Interfaces {
				if _, ok := nameMap[strings.ToLower(name)]; !ok {
					return fmt.Errorf("type '%s' not found in supplied import paths", name)
				}
			}
		}

		if err := generation.Generate(ifaces, opts); err != nil {
			return err
		}
	}

	return nil
}

type loadParams struct {
	fset        *token.FileSet
	importPaths []string

	// gcexportdata specific params
	archives   []archive
	sources    []string
	stdlibRoot string
}

func loadPackages(params loadParams) ([]*internal.GoPackage, error) {
	if len(params.archives) > 0 {
		return PackagesArchive(params)
	}

	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedImports | packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps}, params.importPaths...)
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return nil, errors.New("no packages found")
	}

	ipkgs := make([]*internal.GoPackage, 0, len(pkgs))
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, pkg.Errors[0]
		}
		ipkgs = append(ipkgs, internal.NewPackage(pkg))
	}
	return ipkgs, nil
}
