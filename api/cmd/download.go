package cmd

import (
	"github.com/JSTOR-Labs/pep/api/files"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var dlCmd = &cobra.Command{
	Use:   "download",
	Short: "Download all the necessary files from S3 buckets.",
	Long: `All the files necessary to build the install packages are stored in S3 buckets. 
	This will download them.`,
	Run: func(cmd *cobra.Command, args []string) {
		bucket, err := cmd.Flags().GetString("b")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get private key password")

		}
		err = files.DownloadBucket(bucket, files.Prefix, files.DownloadPath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to download files")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(dlCmd)
	dlCmd.PersistentFlags().String("b", "", "The source bucket for S3 Files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
