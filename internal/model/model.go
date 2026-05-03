package model

import "time"

type ArticleArchive struct {
	Header      string `gorm:"not null"`
	URL         string `gorm:"not null"`
	Text        string
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
type CategoryStat struct {
	Category string `gorm:"not null"`
	Click    uint
	Date     time.Time `gorm:"not null"`
}
type ArticleStat struct {
	URL   string `gorm:"not null"`
	Click uint
	Date  time.Time `gorm:"not null"`
}
