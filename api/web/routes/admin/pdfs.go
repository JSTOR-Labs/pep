package admin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/JSTOR-Labs/pep/api/web/models"
	"github.com/labstack/echo/v4"
)

func CheckPDFs(ctx echo.Context) error {
	req := struct {
		DOIs []string `json:"dois"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	resp := struct {
		DOIs map[string]bool
	}{
		DOIs: make(map[string]bool),
	}

	index, err := pdfs.LoadIndex()
	if err != nil {
		return err
	}

	for _, doi := range req.DOIs {
		resp.DOIs[doi] = index.Search(doi)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func GetPDF(ctx echo.Context) error {
	doi := ctx.Param("doi")
	pdf := ctx.Param("pdf")

	index, err := pdfs.LoadIndex()
	if err != nil {
		return err
	}

	path := index.Get(fmt.Sprintf("%s/%s", doi, pdf))
	if path == "" {
		return ctx.JSON(http.StatusNotFound, models.Response{
			Code:    http.StatusNotFound,
			Message: "PDF not found with that DOI",
		})
	}

	f, err := os.Open(string(path))
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusOK, "application/pdf", data)
}
