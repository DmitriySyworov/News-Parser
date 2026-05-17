package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type ArticleArchive struct {
	Header             string `gorm:"not null"`
	URL                string `gorm:"unique;not null"`
	Text               string
	Category           string    `gorm:"not null"`
	Date               time.Time `gorm:"type:date;not null"`
	ArticleArchiveUUID string    `gorm:"type:char(36);unique;not null"`
	IsArticle          bool
	Error              string `gorm:"-"`
}
type User struct {
	*BaseModel
	Name        string        `gorm:"type:varchar(64);unique;not null"`
	Email       string        `gorm:"type:varchar(256);unique;not null"`
	Password    string        `gorm:"type:char(60);not null"`
	UserUUID    string        `gorm:"type:char(36);primaryKey"`
	UserArticle []UserArticle `gorm:"foreignKey:UserUUID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type UserArticle struct {
	*BaseModel
	Header      string
	URL         string
	Text        string
	Category    string
	ArticleUUID string `gorm:"type:char(36);unique;not null"`
	UserUUID    string `gorm:"type:char(36)"`
	Error       string `gorm:"-"`
}
type CategoryStat struct {
	Category string    `gorm:"not null"`
	Click    uint      `gorm:"type:int"`
	Date     time.Time `gorm:"type:date;not null"`
}
type ArticleStat struct {
	URL   string    `gorm:"not null"`
	Click uint      `gorm:"type:int"`
	Date  time.Time `gorm:"type:date;not null"`
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
	ArticleID uint
	Error     string
}
