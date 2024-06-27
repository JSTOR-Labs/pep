package cmd

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/JSTOR-Labs/pep/api/utils"
	"github.com/JSTOR-Labs/pep/api/web/routes"
	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var incrementHandler = map[string]func(semver.Version) semver.Version{
	"patch": semver.Version.IncPatch,
	"minor": semver.Version.IncMinor,
	"major": semver.Version.IncMajor,
}

// serveCmd represents the serve command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Create and update manifest.json",
	Long:  `If a manifest.json file exists, the script will increment the version number. Otherwise, it will create a new manifest.`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := cmd.Flags().GetString("v")
		if err != nil || v == "" {
			v = "patch"
		}

		hm, err := HasManifest()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to verify manifest file")
		}
		if !hm {
			err := createManifest()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create manifest file")
			}
		} else {
			incrementVersion(v)
		}
	},
}

func HasManifest() (bool, error) {
	hasManifest := false
	path, err := utils.GetManifestPath()
	if err != nil {
		return hasManifest, err
	}

	if _, err := os.Stat(path); err == nil {
		hasManifest = true
	}
	return hasManifest, err
}

func incrementVersion(t string) error {
	path, err := utils.GetManifestPath()
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var m routes.Manifest
	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	m.Version = incrementHandler[t](*semver.MustParse(m.Version)).String()
	m.Updated = time.Now()

	err = writeManifest(m)
	return err
}

func writeManifest(m routes.Manifest) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	path, err := utils.GetManifestPath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	return err
}

func createManifest() error {
	m := routes.Manifest{
		Updated:     time.Now(),
		Name:        "pep",
		Description: "This tool includes an offline browser application that provides a searchable index of content on JSTOR with optional PDF content.",
		Version:     "3.0.1",
		Package:     "github.com/JSTOR-Labs/pep",
	}

	writeManifest(m)

	return nil
}
func init() {
	rootCmd.AddCommand(manifestCmd)
	manifestCmd.PersistentFlags().String("v", "", "Increment type: patch, minor, major")
}
