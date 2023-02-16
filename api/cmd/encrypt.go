package cmd

import (
	"fmt"

	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt all the pdfs using a specified password, then encrypt the password and store it",
	Long: `All PDFs in the PDFs folder will be encrypted using user and owner
	 passwords specified in a flag. The owner password will be forgotten, while
	 the user password will be further encrypted and saved, along with the private RSA
	 key and the Cert used for the encryption.`,
	Run: func(cmd *cobra.Command, args []string) {

		path := "./pdfs"
		err := pdfs.EncryptPDFDirectory(path)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to encrypt PDFs")
		}
		fmt.Println("Document encryption complete")
	},
}

func init() {
	encryptCmd.Flags()
	rootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
