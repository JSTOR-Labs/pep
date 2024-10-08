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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	password   string
	signingKey string
	esAddr     string
	sniff      bool
	webdir     string
	output     string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a config file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("auth.password", password)
		viper.Set("auth.signing_key", signingKey)
		viper.Set("elasticsearch.address", esAddr)
		viper.Set("elasticsearch.sniff", sniff)
		viper.Set("web.root", webdir)
		if output == "" {
			viper.SafeWriteConfig()
		} else {
			viper.WriteConfigAs(output)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&password, "password", "p", "", "Set the admin password")
	generateCmd.Flags().StringVarP(&signingKey, "signing_key", "k", "", "Set the signing key")
	generateCmd.Flags().StringVarP(&esAddr, "elasticsearch_address", "e", "http://localhost:9200", "Set the elasticsearch address")
	generateCmd.Flags().BoolVarP(&sniff, "elasticsearch_sniff", "s", false, "Enable sniff for the elasticsearch connection")
	generateCmd.Flags().StringVarP(&webdir, "webroot", "w", "app", "Set the webroot directory")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "Set the output file")

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
