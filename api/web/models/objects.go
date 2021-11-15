package models

import "github.com/JSTOR-Labs/pep/api/database"

// Response represents a response to a submission endpoint
type (
	Response struct {
		Code    int32       `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	// LoginResponse represents a login response
	LoginResponse struct {
		Token string `json:"token"`
	}

	// Login represents a login request
	Login struct {
		Password string `json:"password"`
	}

	// RequestSubmission represents a students request
	RequestSubmission struct {
		Name     string   `json:"name"`
		Notes    string   `json:"notes"`
		Articles []string `json:"articles"`
	}

	// RequestUpdate represents an admin request to update a request to change its status
	RequestUpdate struct {
		RequestID int                    `json:"requestID"`
		ArticleID string                 `json:"articleID"`
		Status    database.ArticleStatus `json:"status"`
	}

	FormatRequest struct {
		DriveName string `json:"drive"`
	}

	IndexStatusReq struct {
		Name string `json:"name"`
	}

	SyncRequest struct {
		Indices map[string]struct {
			IncludeContent bool `json:"includeContent"`
		} `json:"indices"`
	}
)
