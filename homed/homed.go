package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/JSTOR-Labs/pep/homed/constants"
	"github.com/JSTOR-Labs/pep/homed/netdev"
	"github.com/JSTOR-Labs/pep/homed/payloads"
	"github.com/rs/zerolog/log"
)

func NewDaemon(client *http.Client, done chan bool) *Daemon {
	return &Daemon{client: client, done: done}
}

type Daemon struct {
	client *http.Client
	done   chan bool
}

func (d *Daemon) Run() {
	log.Info().Msg("Starting homed")
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.ping()
		case <-d.done:
			log.Info().Msg("Shutting down homed")
			return
		}
	}
}

func (d *Daemon) checkErrors() bool {
	if f, err := os.Open("/mnt/errors.log"); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}

func (d *Daemon) ping() {
	log.Info().Msg("Sending home ping")

	buf := &bytes.Buffer{}

	mac, err := netdev.GetMACAddress()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get MAC address")
		return
	}
	pingPayload := payloads.Ping{
		MACAddress: mac.String(),
		HasErrors:  d.checkErrors(),
		Version: payloads.VersionInfo{
			// TODO: Get version from somewhere
			Date: time.Now(),
		},
	}

	if err := json.NewEncoder(buf).Encode(pingPayload); err != nil {
		log.Error().Err(err).Msg("Failed to encode ping payload")
		return
	}

	resp, err := d.client.Post(constants.DevicesPingEp, "application/json", buf)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send ping")
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var response payloads.PingResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Error().Err(err).Msg("Failed to decode ping response")
			return
		}
		d.processActions(response.Actions)
	case http.StatusNoContent:
		// Empty response, all is good
	default:
		log.Error().Msgf("Unexpected response code: %d", resp.StatusCode)
	}
}

func (d *Daemon) processActions(actions []payloads.ResponseAction) {
	for _, action := range actions {
		switch action.Op {
		case payloads.OpUpdate:
			var updateData payloads.UpdateData
			if err := json.Unmarshal(action.Data, &updateData); err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal update data")
				continue
			}
			// Fetch the update
			for _, asset := range updateData.Version.Assets {
				if err := asset.Process(d.client); err != nil {
					log.Error().Err(err).Str("asset", asset.FileName).Msg("Failed to process asset")
					continue
				}
			}
		default:
			log.Error().Msgf("Unknown action op: %s", action.Op)
		}
	}
}
