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
	"os"
	"time"

	"github.com/ithaka/labs-pep/updated/storage"
	"github.com/spf13/cobra"
)

// prepareCmd represents the prepare command
var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare's a manifest for a new update",
	Long:  `Prepares a manifest for a new update prior to building a new root`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Create("manifest.json")
		if err != nil {
			return err
		}
		defer f.Close()

		enc := json.NewEncoder(f)

		manifest := &storage.Manifest{
			Updated: time.Now(),
			Version: args[0],
		}

		return enc.Encode(manifest)
	},
}

func init() {
	rootCmd.AddCommand(prepareCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prepareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prepareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
