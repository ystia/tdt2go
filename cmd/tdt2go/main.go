/*
Copyright © 2019 Bull S.A.S. Atos Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ystia/tdt2go"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

var cfgFile string

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var generatedFile string
	var packageName string
	rootCmd = &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "tdt2go <tosca_file>",
		Short: "Generate Go structures from TOSCA datatypes",
		Long:  `tdt2go allows to generate Go source files containing data structures generated from files containing TOSCA data types`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := make([]tdt2go.Option, 0)
			if generatedFile != "" {
				o, err := tdt2go.OutputToFile(generatedFile, 0664)
				if err != nil {
					return err
				}
				opts = append(opts, o)
			}
			if packageName != "" {
				opts = append(opts, tdt2go.Package(packageName))
			}
			return tdt2go.GenerateFile(args[0], opts...)
		},
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tdt2go.yaml)")

	rootCmd.Flags().StringVarP(&generatedFile, "file", "f", "", "file to be generated, if not defined resulting generated file will be printed on default output.")
	rootCmd.Flags().StringVarP(&packageName, "package", "p", "", "package name as it should appear in source file, defaults to the package name of the current directory.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tdt2go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tdt2go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
