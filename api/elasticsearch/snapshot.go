package elasticsearch

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/JSTOR-Labs/pep/api/globals"
	"github.com/JSTOR-Labs/pep/api/utils"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
)

func GetSnapshotStatus(name string) (string, error) {
	resp, err := globals.ES.SnapshotGet("flashdrive").Snapshot(name).Do(context.Background())
	if err != nil {
		return "error", err
	}

	return resp.Snapshots[0].State, nil
}

func CreateSnapshot(indices ...string) (string, error) {
	if len(indices) == 0 {
		return "", errors.New("must specify at least one index")
	}
	_, err := globals.ES.SnapshotGetRepository("flashdrive").Do(context.Background())
	if err != nil {
		// Repo doesn't exist, create it
		_, err := globals.ES.SnapshotCreateRepository("flashdrive").
			Type("fs").Setting("location", "/mnt/es_backup/flashdrive").
			Do(context.Background())
		if err != nil {
			return "", err
		}
	}

	args := make(map[string]string)
	args["indices"] = strings.Join(indices, ",")

	// We now have a valid snapshot repo, time to make the snapshot
	snapshotName := "student-" + utils.RandString(16)
	req := globals.ES.SnapshotCreate("flashdrive", snapshotName)
	resp, err := req.WaitForCompletion(false).BodyJson(args).Do(context.Background())
	if err != nil {
		return "", err
	}
	if !*(resp.Accepted) {
		return "", errors.New("snapshot not accepted")
	}
	log.Println(*(resp.Accepted))
	return snapshotName, nil
}

func LoadSnapshot(addr string, snapshotName string) error {
	sniff := false
	esConf := &config.Config{
		URL:   addr,
		Sniff: &sniff,
	}
	driveEs, err := elastic.NewClientFromConfig(esConf)
	if err != nil {
		return err
	}

	_, err = driveEs.DeleteIndex("_all").Do(context.Background())
	if err != nil {
		return err
	}

	repo, err := driveEs.SnapshotGetRepository("flashdrive").Do(context.Background())
	if err != nil {
		_, err := driveEs.SnapshotCreateRepository("flashdrive").
			Type("fs").Setting("location", "/mnt/es_backup/flashdrive").
			Do(context.Background())
		if err != nil {
			return err
		}
	}

	log.Println(repo)

	req := driveEs.SnapshotRestore("flashdrive", snapshotName).WaitForCompletion(true)

	_, err = req.Do(context.Background())

	return err
}
