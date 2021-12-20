package models

import "time"

type Asset struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	FileName  string `json:"file_name"`
	FileType  string `json:"file_type"`
	Size      uint   `json:"size"`
	Bucket    string `json:"-"`
	Key       string `json:"-"`
	Dst       string `json:"dst"`
	VersionID uint   `json:"-"`

	CreatedAt time.Time `json:"created_at"`
}
