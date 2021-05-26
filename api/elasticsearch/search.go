package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ithaka/labs-pep/api/globals"
	"github.com/olivere/elastic/v7"
)

func BuildFilters(req SearchRequest) []elastic.Query {
	var filters []elastic.Query

	for _, f := range req.Filters {
		filters = append(filters, elastic.NewQueryStringQuery(strings.ReplaceAll(f, "/", "\\/")))
	}
	for _, f := range titleExclude {
		filters = append(filters, elastic.NewQueryStringQuery(fmt.Sprintf("-title:\"%s\"", f)).DefaultOperator("AND"))
	}
	return filters
}

func NewProcessTerms(req SearchRequest) elastic.Query {
	query := elastic.NewBoolQuery()
	var subquery elastic.Query
	must := make([]elastic.Query, 0)
	should := make([]elastic.Query, 0)
	mustNot := make([]elastic.Query, 0)
	if req.Query != "" && req.Query != "*:*" {
		subquery = elastic.NewQueryStringQuery(strings.ReplaceAll(req.Query, "/", "\\/")).
			DefaultOperator("AND").
			FieldWithBoost("title", 3).
			FieldWithBoost("semanticTerms", 3).
			FieldWithBoost("abstract", 2).
			FieldWithBoost("authors", 2).
			Field("disciplines").
			Field("journal").
			Field("publisher").
			Field("ocr").
			Field("doi")
	} else {
		subquery = elastic.NewMatchAllQuery()
	}

	query.Filter(BuildFilters(req)...)
	must = append(must, subquery)

	query.Must(must...)
	query.Should(should...)
	query.MustNot(mustNot...)

	return query
}

func ParseResults(esRes *elastic.SearchResult) *SearchResult {
	res := new(SearchResult)

	res.Total = esRes.TotalHits()
	res.QueryTime = esRes.TookInMillis
	res.Documents = make([]map[string]interface{}, 0)

	for _, hit := range esRes.Hits.Hits {
		doc := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &doc)
		if err != nil {
			os.Exit(3)
		}
		doc["_id"] = doc["doi"]
		doc["_score"] = *hit.Score
		doc["_index"] = hit.Index
		res.Documents = append(res.Documents, doc)
	}

	res.Facets = make(map[string][]map[string]interface{})
	res.Stats = make(map[string]interface{})
	for agg, aggData := range esRes.Aggregations {
		ad := make(map[string]interface{})
		err := json.Unmarshal(aggData, &ad)
		if err != nil {
			os.Exit(3)
		}
		if _, ok := ad["avg"]; ok {
			res.Stats[agg] = aggData
		} else {
			res.Facets[agg] = make([]map[string]interface{}, 0)
			for _, item := range ad["buckets"].([]interface{}) {
				i := item.(map[string]interface{})
				tempAgg := map[string]interface{}{
					"value": i["key"],
					"count": i["doc_count"],
				}
				res.Facets[agg] = append(res.Facets[agg], tempAgg)
			}
		}
	}

	return res
}

func GetESDocument(doi string) (map[string]interface{}, error) {
	query := elastic.NewTermQuery("doi", doi)
	// Query ES for article title and pages
	esDoc, err := globals.ES.Search().Index("_all").Query(query).RestTotalHitsAsInt(true).Do(context.Background())
	if err != nil {
		return nil, err
	}

	if esDoc.Hits.TotalHits.Value == 0 {
		return nil, errors.New("no documents found")
	}

	var doc map[string]interface{}
	err = json.Unmarshal(esDoc.Hits.Hits[0].Source, &doc)

	return doc, err
}
