// Copyright 2018 Bull S.A.S. Atos Technologies - Bull, Rue Jean Jaures, B.P.68, 78340, Les Clayes-sous-Bois, France.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Options are options you are allowed to customize when generating TOSCA datatypes structures
//
// This package uses the functional options pattern https://github.com/tmrts/go-patterns/blob/master/idiom/functional-options.md
// You can't use Options directly but using Option functions.
type Options struct {
	pkg                  string
	output               io.Writer
	generateBuiltinTypes bool
	includePatterns      []string
	excludePatterns      []string
	nameMappings         map[string]string
}

// Option is a function that is allowed to tweak Options
type Option func(*Options)

// GenerateBuiltinTypes option control if TOSCA builtin types should be generated along with
// other datatypes. This option is false by default.
func GenerateBuiltinTypes(p bool) Option {
	return func(o *Options) {
		o.generateBuiltinTypes = p
	}
}

// ExcludePatterns is a set of regexp patterns of data types fully qualified names to exclude.
// Only non-matching datatypes will be transformed.
// Include patterns have the precedence over exclude patterns.
// Default is empty.
func ExcludePatterns(p []string) Option {
	return func(o *Options) {
		o.excludePatterns = p
	}
}

// IncludePatterns is a set of regexp patterns of data types fully qualified names to include.
// Only matching datatypes will be transformed.
// Include patterns have the precedence over exclude patterns.
// Default is empty
func IncludePatterns(p []string) Option {
	return func(o *Options) {
		o.includePatterns = p
	}
}

// Package is the package name as it should appear in source file.
// Defaults to the package name of the current working directory.
func Package(p string) Option {
	return func(o *Options) {
		o.pkg = p
	}
}

// Output is a writer on which generated code will be dumped.
// Defaults to stdout.
func Output(out io.Writer) Option {
	return func(o *Options) {
		o.output = out
	}
}

// NameMappings is map of regular expression and their corresponding remplacements that will be
// applied to TOSCA datatype fully qualified names to transform them into Go struct names.
//
// Defaults to no mappings.
func NameMappings(mappings map[string]string) Option {
	return func(o *Options) {
		o.nameMappings = mappings
	}
}

// OutputToFile is an helper function that allow to dump generated code into a file
//
// See Output
func OutputToFile(outputFile string, perm os.FileMode) (Option, error) {
	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return nil, err
	}
	return Output(f), nil
}

func defaultOptions(toscaFile string) (*Options, error) {
	p, err := getCurrentPackage()
	if err != nil {
		return nil, err
	}
	o := &Options{
		pkg:    p,
		output: os.Stdout,
	}
	return o, nil
}

// GenerateFile generates go code for TOSCA datatypes contains in the given TOSCA definition file.
//
// Generation could be parametrized using Options.
func GenerateFile(toscaFile string, opts ...Option) error {
	options, err := defaultOptions(toscaFile)
	if err != nil {
		return err
	}
	for _, o := range opts {
		o(options)
	}
	p := &parser.Parser{IncludePatterns: options.includePatterns, ExcludePatterns: options.excludePatterns, NameMappings: options.nameMappings}
	dataTypes, err := p.ParseTypes(toscaFile)
	if err != nil {
		return err
	}
	if options.generateBuiltinTypes {
		dataTypes = append(dataTypes, getBuiltinTypes()...)
	}
	f := model.File{
		Package:   options.pkg,
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
	_, err := options.output.Write(content)
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

func getBuiltinTypes() []model.DataType {
	return []model.DataType{
		{
			Name:        "Range",
			FQDTN:       "tosca:range",
			DerivedFrom: "[]uint64",
		},
		{
			Name:        "ScalarUnit",
			FQDTN:       "tosca:scalar-unit",
			DerivedFrom: "string",
		},
		{
			Name:        "ScalarUnitBitRate",
			FQDTN:       "tosca:scalar-unit.bitrate",
			DerivedFrom: "ScalarUnit",
		},
		{
			Name:        "ScalarUnitFrequency",
			FQDTN:       "tosca:scalar-unit.frequency",
			DerivedFrom: "ScalarUnit",
		},
		{
			Name:        "ScalarUnitSize",
			FQDTN:       "tosca:scalar-unit.size",
			DerivedFrom: "ScalarUnit",
		},
		{
			Name:        "ScalarUnitTime",
			FQDTN:       "tosca:scalar-unit.time",
			DerivedFrom: "ScalarUnit",
		},
		{
			Name:        "Version",
			FQDTN:       "tosca:version",
			DerivedFrom: "string",
		},
	}
}
