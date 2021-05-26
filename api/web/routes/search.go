package routes

import (
	"context"
	"net/http"

	"github.com/ithaka/labs-pep/api/elasticsearch"
	"github.com/ithaka/labs-pep/api/globals"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic/v7"
)

func Search(c echo.Context) error {
	var req elasticsearch.SearchRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	var highlightFields []*elastic.HighlighterField
	for _, field := range req.Highlight {
		highlightFields = append(highlightFields, elastic.NewHighlighterField(field))
	}

	// Default limit to 10
	if req.Limit == 0 {
		req.Limit = 10
	}
	doiFound := false
	q := elasticsearch.NewProcessTerms(req)
	for _, f := range req.Fields {
		if f == "doi" {
			doiFound = true
		}
	}
	if !doiFound {
		req.Fields = append(req.Fields, "doi")
	}
	search := globals.ES.Search().
		ErrorTrace(true).
		Index("_all").
		From(req.Offset).Size(req.Limit).
		Query(q).
		Highlight(elastic.NewHighlight().Fields(highlightFields...)).
		RestTotalHitsAsInt(true).FetchSourceContext(elastic.NewFetchSourceContext(true).Include(req.Fields...))
	for _, facet := range req.Facets {
		search = search.Aggregation(facet, elastic.NewTermsAggregation().Field(facet))
	}

	for _, field := range req.Stats {
		search = search.Aggregation(field, elastic.NewStatsAggregation().Field(field))
	}

	res, err := search.Do(context.Background())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err, "resp": res})
	}

	return c.JSON(http.StatusOK, elasticsearch.ParseResults(res))
}
