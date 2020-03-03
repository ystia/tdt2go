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

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ystia/tdt2go"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var generatedFile string
var packageName string
var includePatterns []string
var excludePatterns []string
var nameMappings map[string]string
var generateBuiltinTypes bool

func init() {

	rootCmd = &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "tdt2go <tosca_file>",
		Short: "Generate Go structures from TOSCA datatypes",
		Long:  `tdt2go allows to generate Go source files containing data structures generated from files containing TOSCA data types`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := generateOptions()
			if err != nil {
				return err
			}
			return tdt2go.GenerateFile(args[0], opts...)
		},
	}

	rootCmd.Flags().StringVarP(&generatedFile, "file", "f", "", "file to be generated, if not defined resulting generated file will be printed on default output.")
	rootCmd.Flags().StringVarP(&packageName, "package", "p", "", "package name as it should appear in source file, defaults to the package name of the current directory.")
	rootCmd.Flags().StringSliceVarP(&includePatterns, "include", "i", nil, "regexp patterns of data types fully qualified names to include. Only matching datatypes will be transformed. Include patterns have the precedence over exclude patterns.")
	rootCmd.Flags().StringSliceVarP(&excludePatterns, "exclude", "e", nil, "regexp patterns of data types fully qualified names to exclude. Only non-matching datatypes will be transformed. Include patterns have the precedence over exclude patterns.")
	rootCmd.Flags().BoolVarP(&generateBuiltinTypes, "generate-builtin", "b", false, "Generate tosca builtin types as 'range' or 'scalar-unit' for instance along with datatypes in this file. (default: false)")
	rootCmd.Flags().StringToStringVarP(&nameMappings, "name-mappings", "m", nil, "map of regular expressions and their corresponding remplacements that will be applied to TOSCA datatypes fully qualified names to transform them into Go struct names. This is generally used to keep information from the fully qualified name into the generated name.")
}

func generateOptions() ([]tdt2go.Option, error) {
	opts := make([]tdt2go.Option, 0)
	if generatedFile != "" {
		o, err := tdt2go.OutputToFile(generatedFile, 0664)
		if err != nil {
			return nil, err
		}
		opts = append(opts, o)
	}
	if packageName != "" {
		opts = append(opts, tdt2go.Package(packageName))
	}
	if excludePatterns != nil {
		opts = append(opts, tdt2go.ExcludePatterns(excludePatterns))
	}
	if includePatterns != nil {
		opts = append(opts, tdt2go.IncludePatterns(includePatterns))
	}
	if generateBuiltinTypes {
		opts = append(opts, tdt2go.GenerateBuiltinTypes(true))
	}
	if nameMappings != nil {
		opts = append(opts, tdt2go.NameMappings(nameMappings))
	}
	return opts, nil
}
