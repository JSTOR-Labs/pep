package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/JSTOR-Labs/pep/api/utils"
	"github.com/JSTOR-Labs/pep/api/web/models"

	"github.com/labstack/echo/v4"
)

type (
	Manifest struct {
		Updated     time.Time `json:"updated"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Version     string    `json:"version"`
		Package     string    `json:"package"`
	}

	VersionData struct {
		Version     string    `json:"version"`
		HasContent  bool      `json:"has_content"`
		LastUpdated time.Time `json:"lastUpdate"`
	}
)

func VersionInfo(c echo.Context) error {
	var r Manifest
	path, err := utils.GetManifestPath()
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to open version file",
			Data:    err.Error(),
		})
	}

	data, err := io.ReadAll(f)
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

	hasPDFs, err := pdfs.HasPDFs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check for PDFs",
			Data:    err.Error(),
		})
	}
	v := VersionData{
		Version:     r.Version,
		HasContent:  hasPDFs,
		LastUpdated: r.Updated,
	}

	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    v,
	})
}
