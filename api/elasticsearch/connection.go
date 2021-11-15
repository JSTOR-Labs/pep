package elasticsearch

import (
	"os"
	"time"

	"github.com/JSTOR-Labs/pep/api/globals"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
)

func Connect() (err error) {
	if viper.GetBool("runtime.flash_drive_mode") {
		attempts := 0
		for {
			globals.ES, err = elastic.NewClient(
				elastic.SetURL(viper.GetString("elasticsearch.address")),
				elastic.SetSniff(viper.GetBool("elasticsearch.sniff")),
			)
			if err != nil {
				//e.Logger.Errorf("Waiting for Elasticsearch to start...")
				attempts++
				if attempts == 12 {
					//e.Logger.Error("Elasticsearch took too long to respond, giving up.")
					//e.Logger.Error("You can close everything and try again, if the issue persists," +
					//	" please contact the JSTOR Labs team.")
					time.Sleep(30 * time.Second)
					os.Exit(1)
				}
				time.Sleep(10 * time.Second)
				continue
			}
			break
		}
	} else {
		globals.ES, err = elastic.NewClient(
			elastic.SetURL(viper.GetString("elasticsearch.address")),
			elastic.SetSniff(viper.GetBool("elasticsearch.sniff")),
		)
		if err != nil {
			// e.Logger.Fatal("Unable to create ES client", err)
			return
		}
	}
	return
}
