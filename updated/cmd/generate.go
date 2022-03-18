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
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/JSTOR-Labs/pep/updated/utils"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [key]",
	Short: "Generate a new signing key",
	Long:  `Generates a new ed25519 signing key`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a key file argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		outFile, err := os.Create(args[0])
		if err != nil {
			return err
		}
		defer outFile.Close()

		outPubFile, err := os.Create(fmt.Sprintf("%s.pub", args[0]))
		if err != nil {
			return err
		}

		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}

		err = utils.SaveKey(outFile, priv)
		if err != nil {
			return err
		}

		err = utils.SaveKey(outPubFile, pub)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
