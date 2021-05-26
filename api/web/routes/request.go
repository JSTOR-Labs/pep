package routes

import (
	"net/http"
	"time"

	"github.com/ithaka/labs-pep/api/database"
	"github.com/ithaka/labs-pep/api/web/models"
	"github.com/labstack/echo/v4"
)

func SubmitRequests(c echo.Context) error {
	var s models.RequestSubmission
	if err := c.Bind(&s); err != nil {
		c.Logger().Error("Unable to bind request: ", err)
		return c.JSON(http.StatusBadRequest, models.Response{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if s.Name == "" || len(s.Articles) == 0 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Requests require a name and at least one article",
		})
	}

	request := new(database.Request)
	request.Name = s.Name
	request.Notes = s.Notes
	request.Submitted = time.Now().UTC()

	var documents []database.RequestedDocument

	for _, a := range s.Articles {
		documents = append(documents, database.RequestedDocument{
			Id:     a,
			Status: database.Pending,
		})
	}

	request.Documents = documents

	if err := database.PutRequest(request); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Code: http.StatusInternalServerError, Message: "Unable to write to DB"})
	}

	return c.JSON(http.StatusOK, models.Response{Code: http.StatusOK, Message: "Accepted"})
}
