package generation

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/derision-test/go-mockgen/internal/mockgen/consts"
	"github.com/derision-test/go-mockgen/internal/mockgen/paths"
	"github.com/derision-test/go-mockgen/internal/mockgen/types"
)

type Options struct {
	PackageOptions []PackageOptions
	OutputOptions  OutputOptions
	ContentOptions ContentOptions
}

type PackageOptions struct {
	ImportPaths []string
	Interfaces  []string
	Exclude     []string
	Prefix      string
}

type OutputOptions struct {
	OutputFilename    string
	OutputDir         string
	Force             bool
	DisableFormatting bool
	GoImportsBinary   string
	ForTest           bool
}

type ContentOptions struct {
	PkgName           string
	OutputImportPath  string
	Prefix            string
	ConstructorPrefix string
	FilePrefix        string
}

func Generate(ifaces []*types.Interface, opts *Options) error {
	if opts.OutputOptions.OutputFilename != "" {
		return generateFile(ifaces, opts)
	}

	return generateDirectory(ifaces, opts)
}

func generateFile(ifaces []*types.Interface, opts *Options) error {
	basename := opts.OutputOptions.OutputFilename
	if opts.OutputOptions.ForTest {
		ext := filepath.Ext(basename)
		basename = strings.TrimSuffix(basename, ext) + "_test" + ext
	}

	filename := filepath.Join(opts.OutputOptions.OutputDir, basename)

	exists, err := paths.Exists(filename)
	if err != nil {
		return err
	}
	if exists && !opts.OutputOptions.Force {
		return fmt.Errorf("filename %s already exists, overwrite with --force", paths.GetRelativePath(filename))
	}

	return generateAndRender(ifaces, filename, opts)
}

func generateDirectory(ifaces []*types.Interface, opts *Options) error {
	suffix := "_mock"
	if opts.OutputOptions.ForTest {
		suffix += "_test"
	}

	makeFilename := func(iface *types.Interface) string {
		prefix := opts.ContentOptions.Prefix
		if iface.Prefix != "" {
			prefix = iface.Prefix
		}
		if prefix != "" {
			prefix += "_"
		}

		filename := fmt.Sprintf("%s%s%s.go", prefix, iface.Name, suffix)
		return path.Join(opts.OutputOptions.OutputDir, strings.Replace(strings.ToLower(filename), "-", "_", -1))
	}

	if !opts.OutputOptions.Force {
		allPaths := make([]string, 0, len(ifaces))
		for _, iface := range ifaces {

			allPaths = append(allPaths, makeFilename(iface))
		}

		conflict, err := paths.AnyExists(allPaths)
		if err != nil {
			return err
		}
		if conflict != "" {
			return fmt.Errorf("filename %s already exists, overwrite with --force", paths.GetRelativePath(conflict))
		}
	}

	for _, iface := range ifaces {
		if err := generateAndRender([]*types.Interface{iface}, makeFilename(iface), opts); err != nil {
			return err
		}
	}

	return nil
}

func generateAndRender(ifaces []*types.Interface, filename string, opts *Options) error {
	pkgName := opts.ContentOptions.PkgName
	if opts.OutputOptions.ForTest {
		pkgName += "_test"
	}

	content, err := generateContent(ifaces, pkgName, opts.ContentOptions.Prefix, opts.ContentOptions.ConstructorPrefix, opts.ContentOptions.FilePrefix, opts.ContentOptions.OutputImportPath)
	if err != nil {
		return err
	}

	log.Printf("writing to '%s'\n", paths.GetRelativePath(filename))
	if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}

	if !opts.OutputOptions.DisableFormatting {
		if err := exec.Command(opts.OutputOptions.GoImportsBinary, "-w", filename).Run(); err != nil {
			return errorWithSolutions{
				err: fmt.Errorf("failed to format file: %s", err),
				solutions: []string{
					"install goimports on your PATH",
					"specify a non-standard path to a goimports binary via --goimports",
					"disable post-render formatting via --disable-formatting",
				},
			}
		}
	}

	return nil
}

func generateContent(ifaces []*types.Interface, pkgName, prefix, constructorPrefix, fileContentPrefix, outputImportPath string) (string, error) {
	if fileContentPrefix != "" {
		separator := "\n// "
		fileContentPrefix = "\n//" + separator + strings.Join(strings.Split(strings.TrimSpace(fileContentPrefix), "\n"), separator)
	}

	file := jen.NewFile(pkgName)
	file.HeaderComment(fmt.Sprintf("// Code generated by %s %s; DO NOT EDIT.%s", consts.Name, consts.Version, fileContentPrefix))

	for _, iface := range ifaces {
		log.Printf("generating code for interface '%s'\n", iface.Name)
		generateInterface(file, iface, prefix, constructorPrefix, outputImportPath)
	}

	buffer := &bytes.Buffer{}
	if err := file.Render(buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func generateInterface(file *jen.File, iface *types.Interface, prefix, constructorPrefix, outputImportPath string) {
	if iface.Prefix != "" {
		// Override parent prefix if one is set on the iface
		prefix = iface.Prefix
	}

	withConstructorPrefix := func(f func(*wrappedInterface, string, string) jen.Code) func(*wrappedInterface, string) jen.Code {
		return func(iface *wrappedInterface, outputImportPath string) jen.Code {
			return f(iface, constructorPrefix, outputImportPath)
		}
	}

	topLevelGenerators := []func(*wrappedInterface, string) jen.Code{
		generateMockStruct,
		withConstructorPrefix(generateMockStructConstructor),
		withConstructorPrefix(generateMockStructStrictConstructor),
		withConstructorPrefix(generateMockStructFromConstructor),
	}

	methodGenerators := []func(*wrappedInterface, *wrappedMethod, string) jen.Code{
		generateMockFuncStruct,
		generateMockInterfaceMethod,
		generateMockFuncSetHookMethod,
		generateMockFuncPushHookMethod,
		generateMockFuncSetReturnMethod,
		generateMockFuncPushReturnMethod,
		generateMockFuncNextHookMethod,
		generateMockFuncAppendCallMethod,
		generateMockFuncHistoryMethod,
		generateMockFuncCallStruct,
		generateMockFuncCallArgsMethod,
		generateMockFuncCallResultsMethod,
	}

	titleName := strings.ToUpper(string(iface.Name[0])) + iface.Name[1:]
	mockStructName := fmt.Sprintf("Mock%s%s", prefix, titleName)
	wrappedInterface := wrapInterface(iface, prefix, titleName, mockStructName, outputImportPath)

	for _, generator := range topLevelGenerators {
		file.Add(generator(wrappedInterface, outputImportPath))
		file.Line()
	}

	for _, method := range wrappedInterface.wrappedMethods {
		for _, generator := range methodGenerators {
			file.Add(generator(wrappedInterface, method, outputImportPath))
			file.Line()
		}
	}
}
