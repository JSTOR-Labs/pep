package admin

import (
	"fmt"
	"net/http"

	"github.com/JSTOR-Labs/pep/api/database"
	"github.com/JSTOR-Labs/pep/api/web/models"
	"github.com/labstack/echo/v4"
)

func AdminGetRequests(c echo.Context) error {
	pending := c.QueryParam("pending")

	adminRequests, err := database.GetAdminRequests(pending)

	if err != nil {
		c.Logger().Error("failed to get admin requests: ", err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to fulfill request",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{Code: http.StatusOK, Message: "OK", Data: adminRequests})
}

func isRequestComplete(req *database.Request) bool {
	for _, doc := range req.Documents {
		if doc.Status == database.Pending {
			return false
		}
	}
	return true
}

func AdminUpdateRequest(c echo.Context) error {
	var u models.RequestUpdate
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if u.Status < database.Print || u.Status > database.Denied {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Status must be at least %d and at most %d", database.Print, database.Denied),
		})
	}

	var r database.Request
	err := r.Get(u.RequestID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get request",
			Data:    err.Error(),
		})
	}

	before := isRequestComplete(&r)

	for count, d := range r.Documents {
		if d.Id == u.ArticleID {
			r.Documents[count].Status = u.Status
		}
	}

	after := isRequestComplete(&r)

	if before != after {
		// Remove from pending DB
		database.RemovePendingRequest(&r)
	}

	err = r.Save()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to write request",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{Code: http.StatusOK, Message: "Updated"})
}
