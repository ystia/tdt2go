package tdt2go

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/ystia/tdt2go/internal/pkg/generator"
	"github.com/ystia/tdt2go/internal/pkg/model"
	"github.com/ystia/tdt2go/internal/pkg/parser"
	"golang.org/x/tools/go/packages"
)

type Options struct {
	Package         string
	Output          io.Writer
	IncludePatterns []string
	ExcludePatterns []string
}

type Option func(*Options)

func ExcludePatterns(p []string) Option {
	return func(o *Options) {
		o.ExcludePatterns = p
	}
}

func IncludePatterns(p []string) Option {
	return func(o *Options) {
		o.IncludePatterns = p
	}
}

func Package(p string) Option {
	return func(o *Options) {
		o.Package = p
	}
}

func Output(out io.Writer) Option {
	return func(o *Options) {
		o.Output = out
	}
}

func OutputToFile(outputFile string, perm os.FileMode) (Option, error) {
	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return nil, err
	}
	return func(o *Options) {
		o.Output = f
	}, nil
}

func defaultOptions(toscaFile string) (*Options, error) {
	p, err := getCurrentPackage()
	if err != nil {
		return nil, err
	}
	o := &Options{
		Package: p,
		Output:  os.Stdout,
	}
	return o, nil
}

func GenerateFile(toscaFile string, opts ...Option) error {
	options, err := defaultOptions(toscaFile)
	if err != nil {
		return err
	}
	for _, o := range opts {
		o(options)
	}
	p := &parser.Parser{IncludePatterns: options.IncludePatterns, ExcludePatterns: options.ExcludePatterns}
	dataTypes, err := p.ParseTypes(toscaFile)
	if err != nil {
		return err
	}

	f := generator.File{
		Package:   options.Package,
		Imports:   getImports(dataTypes),
		DataTypes: dataTypes,
	}

	g := &generator.Generator{}
	content, err := g.GenerateFile(f)
	if err != nil {
		return err
	}
	err = outputFile(content, options)
	if err != nil {
		return fmt.Errorf("failed to write generated file: %w", err)
	}
	return nil
}

func outputFile(content []byte, options *Options) error {
	_, err := options.Output.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write generated content: %w", err)
	}
	return nil
}

func getImports(dataTypes []model.DataType) []string {
	imports := make(sort.StringSlice, 0)
	for _, dt := range dataTypes {
		for _, f := range dt.Fields {
			i := getImportForType(f.Type)
			if i != "" && !strSliceContains(imports, i) {
				imports = append(imports, i)
			}
		}
	}
	sort.Sort(imports)
	return imports
}

func getImportForType(t string) string {
	switch {
	case strings.HasPrefix(t, "time."):
		return "time"
	default:
		return ""
	}
}

func strSliceContains(s []string, elem string) bool {
	for _, e := range s {
		if e == elem {
			return true
		}
	}
	return false
}

func getCurrentPackage() (string, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedName,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return "", fmt.Errorf("failed to load current package: %w", err)
	}
	if len(pkgs) != 1 {
		return "", fmt.Errorf("failed to load current package: %d packages found", len(pkgs))
	}
	return pkgs[0].Name, nil
}
