package payloads

import "github.com/JSTOR-Labs/pep/home-api/models"

type ActionOp int

const (
	// Reboot the NUC
	OpReboot ActionOp = iota
	// Fetch a file from an internet source
	OpFetch
	// Verify the checksum of a file
	OpVerify
	// Extract a tar archive
	OpExtract
	// Delete files or directories from disk
	OpDelete
	// Execute a command on the NUC
	OpExec
	// Put a file to an internet destination
	OpPut
	// Perform an Update
	OpUpdate
)

type ResponseAction struct {
	Op   ActionOp    `json:"op"`
	Data interface{} `json:"data,omitempty"`
}

type PingResponse struct {
	Actions []ResponseAction `json:"actions,omitempty"`
}

type FetchData struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type VerifyData struct {
	Path     string `json:"path"`
	Checksum []byte `json:"checksum"`
}

type ExtractData struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type DeleteData struct {
	Path string `json:"path"`
}

type ExecData struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
}

type PutData struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type UpdateData struct {
	models.Version
}
