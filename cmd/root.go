/*
Copyright © 2024 Jason Jackson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jason-jackson/makego/src"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string

var project = src.NewProject()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "makego [flags] [package_name]",
	Short: "A customizable code generator to quickly set up APIs in Go.",
	Long: `Makego is a customizable code generator that sets up the basics of an API framework in Go
to let you quickly launch an api.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Return values from viper back to cobra if needed
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if !f.Changed && viper.IsSet(f.Name) {
				f.Value.Set(viper.GetString(f.Name))
			}
		})

		return viper.UnmarshalKey("templates", &project.Templates)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			project.PkgName = args[0]
		}

		return project.Generate()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Expect either 0 (allowed if `go.mod` already exists) or 1, for package name
	cobra.MaximumNArgs(1)

	// Set all flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/makego.yaml)")

	rootCmd.Flags().StringVar(&project.AppName, "name", "", "application name")
	rootCmd.Flags().StringVar(&project.Folder, "folder", "", "application folder, can be left blank for no folder")
	rootCmd.Flags().StringVar(&project.License, "license", "", "license, can be left blank for proprietary code")
	rootCmd.Flags().StringVar(&project.Copyright, "copyright", "", "copyright holder (and contact if desired)")
	rootCmd.Flags().StringVar(&project.Database.Name, "database", "postgres", "database type to use (mysql, mariadb, postgres, etc)")
	rootCmd.Flags().StringVar(&project.ORM.Name, "orm", "gorm", "ORM to use for models (defaults to gorm)")
	rootCmd.Flags().StringVar(&project.Router.Name, "router", "gin", "router to use (echo, gin, http, mux)")
	rootCmd.Flags().StringVar(&project.EnvPrefix, "envprefix", "", "how to expect env variables to be prefixed")

	rootCmd.Flags().BoolVarP(&project.Docker, "docker", "d", false, "whether to use docker")
	rootCmd.Flags().BoolVarP(&project.Header, "header", "a", false, "whether to show copyright headers on most files")
	rootCmd.Flags().BoolVarP(&project.Sentry, "sentry", "s", false, "whether to use sentry")

	err := viper.BindPFlags(rootCmd.Flags())
	cobra.CheckErr(err)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".makego" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("makego")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
