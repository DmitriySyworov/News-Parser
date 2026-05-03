package stat

import (
	"time"
)

type ResponseStatCategoryDate struct {
	Categories []CategoryDbDate
	Error      string
}
type ResponseStatCategoryAll struct {
	Categories []CategoryDbAll
	Error      string
}
type CategoryDbDate struct {
	Category string
	Click    uint
	Date     time.Time
	Place    uint
}
type CategoryDbAll struct {
	Category string
	SumClick uint
	Place    uint
}
type ResponseStatArticleDate struct {
	Articles []ArticleDbDate
	Error    string
}
type ResponseStatArticleAll struct {
	Articles []ArticleDbAll
	Error    string
}
type ArticleDbDate struct {
	URL   string
	Click uint
	Date  time.Time
	Place uint
}
type ArticleDbAll struct {
	URL      string
	SumClick uint
	Place    uint
}
