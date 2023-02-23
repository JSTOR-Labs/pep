package routes

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/JSTOR-Labs/pep/api/web/models"
	"github.com/labstack/echo/v4"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

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

	pw, err := pdfs.GetPDFPassword(nil)
	if err != nil {
		return err
	}

	rs := pdfs.GetReadSeeker(f)
	var buf bytes.Buffer

	config := pdfcpu.LoadConfiguration()
	config.UserPW = pw

	err = pdfcpu.Decrypt(rs, &buf, config)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, "application/pdf", buf.Bytes())
}
