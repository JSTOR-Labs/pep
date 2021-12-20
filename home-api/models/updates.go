package models

import "time"

type Version struct {
	ID          uint      `gorm:"primaryKey"`
	Version     string    `json:"version"`
	Date        time.Time `json:"date"`
	Assets      []Asset   `json:"assets"`
	DependsOnID *uint     `json:"-"`
	DependsOn   *Version  `gorm:"foreignKey:DependsOnID" json:"-"`
}
