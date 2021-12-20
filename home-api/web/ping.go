package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/JSTOR-Labs/pep/home-api/models"
	"github.com/JSTOR-Labs/pep/home-api/payloads"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (s *Web) PingHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var ping payloads.Ping

	err := dec.Decode(&ping)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device := models.Device{}
	if err := s.First(&device, "mac_address = ?", ping.MACAddress).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		device.MACAddress = ping.MACAddress
		device.LastSeen = time.Now()
		device.Institution = "Unk"
		device.Version = models.VersionInfo{
			Version: ping.Version.Version,
			Date:    ping.Version.Date,
		}
		if err := s.Create(&device).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		device.LastSeen = time.Now()
		device.Version = models.VersionInfo{
			Version: ping.Version.Version,
			Date:    ping.Version.Date,
		}
		if err := s.Save(&device).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	update := models.Version{}
	if err := s.Where("date > ?", ping.Version.Date).Order("date desc").Limit(1).First(&update).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("failed to get updates")
			// Fail gracefully, NUC doesn't need to know that the server-side failed here
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// We have the updates, let's create a response payload
	// Check latest update for deps
	if update.DependsOnID != nil {
		// Update has deps, let's get the earliest dep that's after the current version
		for {
			nextUpdate := models.Version{}
			if err := s.First(&nextUpdate, update.DependsOnID).Error; err != nil {
				log.Error().Err(err).Msg("failed to get updates")
				// Fail gracefully, NUC doesn't need to know that the server-side failed here
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if nextUpdate.Date.Before(update.Date) {
				break
			}
			update = nextUpdate
		}
	}

	actions := []payloads.ResponseAction{
		{
			Op: payloads.OpUpdate,
			Data: payloads.UpdateData{
				Version: update,
			},
		},
	}
	payload := payloads.PingResponse{
		Actions: actions,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
	w.WriteHeader(http.StatusOK)
}
