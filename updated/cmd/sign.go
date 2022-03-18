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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/JSTOR-Labs/pep/updated/utils"
	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign [file] [key] [outdir]",
	Short: "Sign a file using the specified private key",
	Long:  `Signs a file using the specified ed25519 private key`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		input, key := args[0], args[1]
		inFile, err := os.Open(input)
		if err != nil {
			return err
		}
		defer inFile.Close()
		keyFile, err := os.Open(key)
		if err != nil {
			return err
		}
		defer keyFile.Close()
		sigFile, err := os.Create(fmt.Sprintf("%s.sig", input))
		if err != nil {
			return err
		}
		defer sigFile.Close()

		privateKey, err := utils.LoadPrivKey(keyFile)
		if err != nil {
			return err
		}

		fileData, err := ioutil.ReadAll(inFile)
		if err != nil {
			return err
		}

		sigData := ed25519.Sign(privateKey, fileData)

		_, err = sigFile.Write(sigData)
		return err
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
