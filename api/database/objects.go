package database

import (
	"time"

	"github.com/boltdb/bolt"
)

type ArticleStatus int32

const ( // Approved is < 2
	Print ArticleStatus = iota
	PDF
	Pending
	Denied
)

const (
	PendingBucket   = "PendingRequests"
	CompletedBucket = "CompletedRequests"
)

type (
	Request struct {
		Id        int
		Name      string
		Notes     string
		Submitted time.Time
		Documents []RequestedDocument
	}

	RequestedDocument struct {
		Id     string
		Status ArticleStatus
	}

	// AdminRequest represents a single request that needs to be approved by an admin
	AdminRequest struct {
		ID            string        `json:"id"`
		RequestID     int           `json:"requestID"`
		Title         string        `json:"title"`
		StudentName   string        `json:"studentName"`
		DateRequested time.Time     `json:"dateRequested"`
		NumPages      int           `json:"numPages"`
		Notes         string        `json:"notes"`
		PDFAvailable  bool          `json:"pdf"`
		Status        ArticleStatus `json:"status"`
		SrcHTML       string        `json:"srcHtml"`
	}

	// AdminRequests represents the list of requests an admin needs to approve
	AdminRequests struct {
		Requests []*AdminRequest `json:"requests"`
	}
)

var bdb *bolt.DB
