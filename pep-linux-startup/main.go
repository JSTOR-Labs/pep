package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func getUserID(username string) (int, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return 0, err
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func getGroupID(group string) (int, error) {
	g, err := user.LookupGroup(group)
	if err != nil {
		return 0, err
	}
	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return 0, err
	}
	return gid, nil
}

func recursiveChown(path string, uid, gid int) error {
	err := os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			err = recursiveChown(filepath.Join(path, f.Name()), uid, gid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	// Look up elasticsearch user
	uid, err := getUserID("elasticsearch")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get user id for elasticsearch")
	}

	// Look up elasticsearch group
	gid, err := getGroupID("elasticsearch")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get group id for elasticsearch")
	}

	// Look up jstor user
	juid, err := getUserID("jstor")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get user id for jstor")
	}

	// Look up jstor group
	jgid, err := getGroupID("jstor")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get group id for jstor")
	}

	// Look up elasticsearch data directory
	f, err := os.Open("/opt/elasticsearch/config/elasticsearch.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open elasticsearch.yml")
	}
	defer f.Close()

	esConfig := make(map[string]interface{})
	err = yaml.NewDecoder(f).Decode(&esConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decode elasticsearch.yml")
	}

	dataDir, ok := esConfig["path.data"].(string)
	if !ok {
		log.Fatal().Msg("Failed to get path.data from elasticsearch.yml")
	}

	if !dirExists("/mnt/es_backup") {
		if err := os.Mkdir("/mnt/es_backup", 0755); err != nil {
			log.Fatal().Err(err).Msg("Failed to create /mnt/es_backup")
		}
	}

	if !dirExists("/mnt/data") {
		if err := os.Mkdir("/mnt/es_backup", 0755); err != nil {
			log.Fatal().Err(err).Msg("Failed to create /mnt/data")
		}
	}

	if !dirExists(dataDir) {
		if err := os.Mkdir(dataDir, 0755); err != nil {
			log.Fatal().Err(err).Msg("Failed to create data directory")
		}
	}

	err = recursiveChown(dataDir, uid, gid)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to chown elasticsearch data directory")
	}
	recursiveChown("/mnt/es_backup", uid, gid)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to chown elasticsearch backup directory")
	}
	recursiveChown("/mnt/data", juid, jgid)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to chown pepapi data directory")
	}
}

func dirExists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}
