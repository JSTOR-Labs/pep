package database

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"time"

	"github.com/JSTOR-Labs/pep/api/elasticsearch"
	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	log.Info().Msg("Running Requests Init")
	bdb = setupBDB()
}

func setupBDB() *bolt.DB {
	var err error
	requestsLocation := "requests.db"
	if !viper.GetBool("runtime.flash_drive_mode") {
		requestsLocation = "/mnt/data/" + requestsLocation
	}
	bdb, err = bolt.Open(requestsLocation, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Warn().Err(err).Msg("Unable to open requests database")
	}
	return bdb
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func PutRequest(request *Request) error {
	if bdb == nil {
		bdb = setupBDB()
	}

	err := bdb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(PendingBucket))
		if err != nil {
			return err
		}
		id, _ := b.NextSequence()
		request.Id = int(id)

		buf, err := json.Marshal(request)
		if err != nil {
			return err
		}
		return b.Put(itob(request.Id), buf)
	})

	if err != nil {
		return err
	}

	return bdb.Sync()
}

func GetAdminRequests(pending string) (adminRequests AdminRequests, err error) {
	if bdb == nil {
		bdb = setupBDB()
	}

	adminRequests.Requests = make([]*AdminRequest, 0)

	noIndex := false
	index, err := pdfs.LoadIndex()
	if err != nil {
		noIndex = true
	}

	err = bdb.Update(func(tx *bolt.Tx) error {
		buckets := make([]*bolt.Bucket, 0)
		bucket, err := tx.CreateBucketIfNotExists([]byte(PendingBucket))
		if err != nil {
			return err
		}
		buckets = append(buckets, bucket)
		if pending != "true" {
			bucket, err := tx.CreateBucketIfNotExists([]byte(CompletedBucket))
			if err != nil {
				return err
			}
			buckets = append(buckets, bucket)
		}
		for _, b := range buckets {
			_ = b.ForEach(func(k, v []byte) error {
				var r Request
				_ = json.Unmarshal(v, &r)
				for _, rd := range r.Documents {
					if (pending == "true" && rd.Status != Pending) || (pending != "true" && rd.Status == Pending) {
						continue
					}
					doc, err := elasticsearch.GetESDocument(rd.Id)
					if err != nil {
						// c.Logger().Error("Unable to get document info ", rd.Id, ": ", err)
						continue
					}

					// Handle missing titles
					var title string
					if _, ok := doc["title"]; !ok {
						// Document is likely a book review, just assign a blank title
						title = ""
					} else {
						title, ok = doc["title"].(string)
						if !ok {
							// ??? huh
							title = ""
						}
					}

					// Convert srcHtml
					var srcHtml string = ""
					if src, ok := doc["srcHtml"]; ok {
						if srcString, ok := src.(string); ok {
							srcHtml = srcString
						}
					}

					PDFAvailable := false
					if !noIndex {
						PDFAvailable = index.Search(rd.Id)
					}

					req := &AdminRequest{
						ID:            rd.Id,
						RequestID:     r.Id,
						Title:         title,
						StudentName:   r.Name,
						DateRequested: r.Submitted,
						NumPages:      int(doc["pageCount"].(float64)),
						Notes:         r.Notes,
						PDFAvailable:  PDFAvailable,
						Status:        rd.Status,
						SrcHTML:       srcHtml,
					}
					adminRequests.Requests = append(adminRequests.Requests, req)
				}
				return nil
			})
		}
		return nil
	})

	return
}

func RemovePendingRequest(r *Request) {
	if bdb == nil {
		bdb = setupBDB()
	}

	_ = bdb.Update(func(tx *bolt.Tx) error {
		pb := tx.Bucket([]byte(PendingBucket))
		return pb.Delete(itob(r.Id))
	})
}

func (req *Request) Save() error {
	if bdb == nil {
		bdb = setupBDB()
	}

	bucketName := []byte(CompletedBucket)
	for _, doc := range req.Documents {
		if doc.Status == Pending {
			bucketName = []byte(PendingBucket)
		}
	}
	err := bdb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(req)
		if err != nil {
			return err
		}
		return b.Put(itob(req.Id), encoded)
	})
	return err
}

func (req *Request) Get(id int) error {
	if bdb == nil {
		bdb = setupBDB()
	}

	err := bdb.View(func(tx *bolt.Tx) error {
		pending := tx.Bucket([]byte(PendingBucket))
		completed := tx.Bucket([]byte(CompletedBucket))
		r := pending.Get(itob(id))
		if r != nil {
			err := json.Unmarshal(r, req)
			return err
		}
		r = completed.Get(itob(id))
		if r != nil {
			err := json.Unmarshal(r, req)
			return err
		}
		return errors.New("no such request")
	})
	return err
}
