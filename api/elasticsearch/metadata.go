package elasticsearch

import (
	"context"

	"github.com/JSTOR-Labs/pep/api/globals"
	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/JSTOR-Labs/pep/api/utils"
)

func GetIndexMetadata() ([]*IndexMetadata, error) {
	sizes, err := pdfs.GetPDFSizes()
	if err != nil {
		sizes = make(map[string]int64)
	}

	stats, err := globals.ES.IndexStats().Metric("indexing", "store").Do(context.Background())
	if err != nil {
		return nil, err
	}

	metadata := make([]*IndexMetadata, 0)

	for k, v := range stats.Indices {
		contentSize, ok := sizes[k]
		if !ok {
			contentSize = 0
		}
		metadata = append(metadata, &IndexMetadata{
			Name:        k,
			Size:        float32(utils.ConvertBytes(v.Primaries.Store.SizeInBytes, utils.GIGABYTE)),
			SizeContent: float32(utils.ConvertBytes(contentSize, utils.GIGABYTE)),
		})
	}

	return metadata, nil
}
