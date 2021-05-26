package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func main() {
	if !dirExists("/mnt/es_backup") {
		if err := os.Mkdir("/mnt/es_backup", 0755); err != nil {
			panic("failed to create es_backup dir")
		}
	}

	// Look up elasticsearch user
	u, err := user.Lookup("elasticsearch")
	if err != nil {
		panic("elasticsearch user not found")
	}

	g, err := user.LookupGroup("elasticsearch")
	if err != nil {
		panic("elasticsearch group not found")
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		panic("user ID is not a number")
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		panic("group ID is not a number")
	}

	err = filepath.Walk("/mnt/es_backup", func(path string, info os.FileInfo, err error) error {
		return os.Chown(path, uid, gid)
	})

	if err != nil {
		panic(fmt.Sprintf("walk failed: %v", err))
	}
}

func dirExists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}
