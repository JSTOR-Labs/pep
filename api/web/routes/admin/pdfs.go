package admin

import (
	"net/http"

	"github.com/JSTOR-Labs/pep/api/pdfs"
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
