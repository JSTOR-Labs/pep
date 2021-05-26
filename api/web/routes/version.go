package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ithaka/labs-pep/api/web/models"
	"github.com/labstack/echo/v4"
)

type (
	Manifest struct {
		Updated time.Time
		Version string
		Package string
	}

	VersionData struct {
		Version     string    `json:"version"`
		LastUpdated time.Time `json:"lastUpdate"`
	}
)

func VersionInfo(c echo.Context) error {
	var r Manifest
	f, err := os.Open("/manifest.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to open version file",
			Data:    err.Error(),
		})
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to read version file",
			Data:    err.Error(),
		})
	}

	err = json.Unmarshal(data, &r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Bad version file structure",
			Data:    err.Error(),
		})
	}

	v := VersionData{
		Version:     r.Version,
		LastUpdated: r.Updated,
	}

	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    v,
	})
}
