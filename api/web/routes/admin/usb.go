package admin

import (
	"net/http"

	"github.com/JSTOR-Labs/pep/api/elasticsearch"
	"github.com/JSTOR-Labs/pep/api/globals"
	"github.com/JSTOR-Labs/pep/api/web/models"
	"github.com/labstack/echo/v4"
)

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
