/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ithaka/labs-pep/updated/storage"
	"github.com/ithaka/labs-pep/updated/utils"
	"github.com/spf13/cobra"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package a new release",
	Long:  `Package takes a directory and creates a new release package`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("1 arg expected, %d provided", len(args))
		}
		fi, err := os.Stat(args[0])
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			return fmt.Errorf("%s is not a directory", args[0])
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]

		if _, err := os.Stat("dist"); os.IsNotExist(err) {
			os.Mkdir("dist", os.FileMode(0755))
		}

		filename := fmt.Sprintf("%s.tar.gz", time.Now().Format("2006-02-01_15-04-05"))

		f, err := os.Create(fmt.Sprintf("dist/%s", filename))
		if err != nil {
			return err
		}
		defer f.Close()

		om, err := os.Open("manifest.json")
		if err != nil {
			return err
		}
		defer om.Close()

		dec := json.NewDecoder(om)

		m, err := os.Create("dist/manifest.json")
		if err != nil {
			return err
		}
		defer m.Close()

		enc := json.NewEncoder(m)

		manifest := storage.Manifest{}
		err = dec.Decode(&manifest)
		if err != nil {
			return err
		}

		manifest.Package = filename

		err = enc.Encode(&manifest)
		if err != nil {
			return err
		}

		os.Chdir(src)
		return utils.Tar(".", f)
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
