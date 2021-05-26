package daemon

import (
	"log"
	"time"

	"github.com/ithaka/labs-pep/updated/storage"
	"github.com/ithaka/labs-pep/updated/utils"
)

// Run starts the update daemon
func Run() {
	m, err := storage.NewManager("org.jstor.labs.pep")
	if err != nil {
		panic("failed to create S3 client")
	}
	for {
		time.Sleep(time.Minute)
		log.Println("Checking for updates")
		available, manifest, err := m.Check()
		if err != nil {
			log.Printf("[WARN] Failed to check for update: %v\n", err)
			continue
		}
		if available {
			log.Println("New update available, processing...")
			err := m.Download(manifest)
			if err != nil {
				log.Printf("[WARN] Failed to download latest update: %v\n", err)
				continue
			}
			err = m.Install(manifest)
			if err != nil {
				log.Printf("[WARN] Failed to install latest update: %v\n", err)
				continue
			}
			log.Println("Update successful, rebooting...")
			utils.Reboot()
		}
	}
}
