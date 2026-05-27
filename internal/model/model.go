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
	Category           string    `gorm:"type:varchar(20);not null"`
	Date               time.Time `gorm:"type:date;not null"`
	ArticleArchiveUUID string    `gorm:"type:char(36);unique;not null"`
	IsArticle          bool
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
	Category    string `gorm:"type:varchar(20)"`
	ArticleUUID string `gorm:"type:char(36);unique;not null"`
	UserUUID    string `gorm:"type:char(36)"`
}
type UserArticleStat struct {
	Date        time.Time `gorm:"type:date;uniqueIndex:date_user_uuid"`
	Created     int
	Updated     int
	SoftDeleted int
	HardDeleted int
	Recovered   int
	UserUUID    string `gorm:"type:char(36);uniqueIndex:date_user_uuid"`
}
type CategoryStat struct {
	Category string    `gorm:"type:varchar(20);not null"`
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
}
