package main

import (
	"io/ioutil"
	"os"
	"syscall"
)

// Returns true if the given directory exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Returns the number of directories or files in the given path
func countFiles(path string) int {
	count := 0
	if pathExists(path) {
		files, _ := ioutil.ReadDir(path)
		count = len(files)
	}
	return count
}

func migrateFromArch() {
	// Mount the old upper filesystem
	if !pathExists("/mnt/oldroot") {
		os.Mkdir("/mnt/oldroot", 0755)
	}
	if err := syscall.Mount("/dev/sda2", "/mnt/oldroot", "ext4", syscall.MS_NOATIME, ""); err != nil {
		return
	}

	// Migrate Elasticsearch data
	if countFiles("/mnt/oldroot/upper/var/lib/elasticsearch") > 0 {
		// Elasticsearch data in old location, move to /mnt/elasticsearch
		os.Rename("/mnt/oldroot/upper/var/lib/elasticsearch", "/mnt/elasticsearch")
	}
	// Migrate Requests database
	if pathExists("/mnt/oldroot/upper/home/jstor/requests.db") {
		// Requests database in old location, move to /mnt/requests.db
		os.Rename("/mnt/oldroot/upper/home/jstor/pep/requests.db", "/mnt/requests.db")
	}
}
