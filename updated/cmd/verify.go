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

	"github.com/ithaka/labs-pep/updated/utils"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify [file] [key]",
	Short: "Verify the signature of a file",
	Long:  `Verify the signature of a file using an ed25519 public key`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		input, key := args[0], args[1]
		inFile, err := os.Open(input)
		if err != nil {
			return err
		}
		defer inFile.Close()

		fileData, err := ioutil.ReadAll(inFile)
		if err != nil {
			return err
		}

		sigFile, err := os.Open(fmt.Sprintf("%s.sig", input))
		if err != nil {
			return err
		}
		defer sigFile.Close()

		sigData, err := ioutil.ReadAll(sigFile)
		if err != nil {
			return err
		}

		keyFile, err := os.Open(key)
		if err != nil {
			return err
		}
		defer keyFile.Close()

		pubKey, err := utils.LoadPubKey(keyFile)
		if err != nil {
			return err
		}

		valid := ed25519.Verify(pubKey, fileData, sigData)
		if !valid {
			fmt.Println("invalid signature")
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
