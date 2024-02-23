/*
Copyright © 2024 Aditya Wardianto <aditya.wardianto11@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ditwrd/setunggal/internal/ansible/data"

	"github.com/kluctl/go-embed-python/embed_util"
	"github.com/kluctl/go-embed-python/python"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "setunggal",
	Short: "Setunggal CLI",
	Run: func(cmd *cobra.Command, args []string) {
		tmpDir := filepath.Join(os.TempDir(), "setunggal-embedded")

		p, _ := python.NewEmbeddedPythonWithTmpDir(tmpDir+"-python", true)
		a, _ := embed_util.NewEmbeddedFilesWithTmpDir(data.Data, tmpDir+"-ansible", true)
		p.AddPythonPath(a.GetExtractedPath())

		fmt.Println("----------")
		fmt.Println("Ansible Version Check")
		fmt.Println("----------")
		cmdd := p.PythonCmd("-m", "ansible", "playbook", "--version")
		cmdd.Stdout = os.Stdout
		cmdd.Stderr = os.Stderr
		cmdd.Run()
		fmt.Println()

		cmdList := [][]string{
			{"list"},
			{"install", "community.docker"},
			{"list"},
		}
		for _, cmdItem := range cmdList {
			defaultCmd := []string{"-m", "ansible", "galaxy", "collection"}
			cmdToRun := append(defaultCmd, cmdItem...)

			fmt.Println("----------")
			fmt.Println("Running python", strings.Join(cmdToRun, " "))
			fmt.Println("----------")

			cmdd := p.PythonCmd(cmdToRun...)
			cmdd.Stdout = os.Stdout
			cmdd.Stderr = os.Stderr
			cmdd.Run()

			fmt.Println()
		}
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.setunggal.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".setunggal" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".setunggal")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
