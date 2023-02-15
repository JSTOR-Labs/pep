package cmd

import (
	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var keyCmd = &cobra.Command{
	Use:   "keys",
	Short: "Generate cert, with public and private keys, and encrypted passwords",
	Long: `All PDFs in the PDFs folder will be encrypted using user and owner
	 passwords specified in a flag. The owner password will be forgotten, while
	 the user password will be further encrypted and saved, along with the private RSA
	 key and the Cert used for the encryption.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pdfs.SaveEncryptionFiles()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to handle encryption files")
		}

	},
}

func init() {
	keyCmd.Flags()
	rootCmd.AddCommand(keyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
