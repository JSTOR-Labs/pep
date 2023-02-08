package cmd

// serveCmd represents the serve command
// var encryptCmd = &cobra.Command{
// 	Use:   "encrypt",
// 	Short: "Encrypt all the pdfs using a specified password, then encrypt the password and store it",
// 	Long: `All PDFs in the PDFs folder will be encrypted using user and owner
// 	 passwords specified in a flag. The owner password will be forgotten, while
// 	 the user password will be further encrypted and saved, along with the private RSA
// 	 key and the Cert used for the encryption.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if !viper.GetBool("runtime.flash_drive_mode") {
// 			home, err := homedir.Dir()
// 			cobra.CheckErr(err)
// 			os.Chdir(home)
// 		}
// 		web.Listen(Port)
// 	},
// }

// func init() {
// 	serveCmd.Flags().IntVarP(&Port, "port", "p", 1323, "Port to listen on")
// 	rootCmd.AddCommand(serveCmd)

// 	// Here you will define your flags and configuration settings.

// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
