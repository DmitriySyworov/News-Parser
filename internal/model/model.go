package model

import "time"

type ArticleArchive struct {
	Header      string    `gorm:"not null"`
	URL         string    `gorm:"not null"`
	Text        string    `gorm:"not null"`
	Category    string    `gorm:"not null"`
	Date        time.Time `gorm:"not null"`
	UUIDArticle string    `gorm:"unique,not null"`
	Error       string    `gorm:"-"`
}
type ArticleToday struct {
	Header    string
	URL       string
	Text      string
	Category  string
	IDArticle uint
	Error     string
}
