package cmd

import (
	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/JSTOR-Labs/pep/api/utils"
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
		pw, err := cmd.Flags().GetString("pw")
		if err != nil {
			pw, err = pdfs.PromptUser(false)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get private key password")
			}
		}

		path, err := utils.GetPDFPath()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get pdf path")
		}
		err = pdfs.EncryptPDFDirectory(path, pw)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to encrypt PDFs")
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.PersistentFlags().String("pw", "", "A password to encrypt the private key")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
