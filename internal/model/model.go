package model

import (
	"time"

	"gorm.io/gorm"
)

type ArticleArchive struct {
	Header      string `gorm:"not null"`
	URL         string `gorm:"unique,not null"`
	Text        string
	Category    string    `gorm:"not null"`
	Date        time.Time `gorm:"not null"`
	UUIDArticle string    `gorm:"unique,not null"`
	IsArticle   bool
	Error       string `gorm:"-"`
}
type User struct {
	*gorm.Model
	Name        string        `gorm:"unique;not null"`
	Email       string        `gorm:"unique;not null"`
	Password    string        `gorm:"not null"`
	UUIDUser    string        `gorm:"unique;not null;primaryKey"`
	UserArticle []UserArticle `gorm:"foreignKey:UUIDUser;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
type UserArticle struct {
	*gorm.Model
	Header    string
	URL       string
	Text      string
	Category  string
	IDArticle uint   `gorm:"type:text;unique;not null"`
	UUIDUser  string `gorm:"not null"`
	Error     string `gorm:"-"`
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
type TemporaryData struct {
	Name      string
	Email     string
	Password  string
	TempCode  uint
	IDSession string
}
type ArticleToday struct {
	Header    string
	URL       string
	Text      string
	Category  string
	IDArticle uint
	Error     string
}
