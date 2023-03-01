package cmd

import (
	"github.com/JSTOR-Labs/pep/api/files"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Gather all the necessary files from S3 buckets and assemble download packages",
	Long: `All the files necessary to build the download packages are stored in S3 buckets. 
	This will gather them together and build the individual packages, then upload the zipped
	files to the appropriate bucket.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := files.DownloadBucket(files.Bucket, files.Prefix, files.DownloadPath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to download files")
			return
		}
		err = files.AssembleChromebook(files.DownloadPath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to assemble chromebook files")
			return
		}
	},
}

func init() {
	gatherCmd.Flags()
	rootCmd.AddCommand(gatherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
