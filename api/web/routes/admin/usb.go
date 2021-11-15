package admin

import (
	"log"
	"net/http"

	"github.com/JSTOR-Labs/pep/api/elasticsearch"
	"github.com/JSTOR-Labs/pep/api/globals"
	"github.com/JSTOR-Labs/pep/api/usb"
	"github.com/JSTOR-Labs/pep/api/web/models"
	"github.com/labstack/echo/v4"
)

func GetUSBDevices(c echo.Context) error {
	return c.JSON(http.StatusOK, models.Response{Code: http.StatusOK, Data: usb.FindUSBDrives()})
}

func FormatUSBDevice(c echo.Context) error {
	var req models.FormatRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	err := usb.FormatDrive(req.DriveName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "Done",
		Data:    nil,
	})
}

func GetIndexData(c echo.Context) error {
	data, err := elasticsearch.GetIndexMetadata()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func SnapshotStatus(c echo.Context) error {
	var req models.IndexStatusReq
	if err := c.Bind(&req); err != nil {
		return err
	}
	status, err := elasticsearch.GetSnapshotStatus(req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    map[string]string{"status": status},
	})
}

func GetRestoreStatus(c echo.Context) error {
	snapshotName := c.QueryParam("snapshot")
	val, ok := globals.BuildJobs[snapshotName]
	if !ok {
		return c.JSON(http.StatusOK, models.Response{
			Code: http.StatusOK,
			Data: map[string]bool{
				"done":   false,
				"failed": false,
			},
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Code: http.StatusOK,
		Data: map[string]bool{
			"done":   val == 1,
			"failed": val == 2,
		},
	})
}

func BuildFlashDrive(c echo.Context) error {
	name := c.Param("name")
	var req models.SyncRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	var snapshotName string

	indices := make([]string, 0)
	pdfs := make([]string, 0)

	for k, v := range req.Indices {
		indices = append(indices, k)
		if v.IncludeContent {
			pdfs = append(pdfs, k)
		}
	}

	if snapshot, err := elasticsearch.CreateSnapshot(indices...); err != nil {
		log.Println("failed to create initial snapsnot: ", err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
	} else {
		snapshotName = snapshot
	}

	err := usb.BuildFlashDrive(name, snapshotName, pdfs)
	if err != nil {
		log.Println("error initializing drive: ", err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    snapshotName,
	})
}
