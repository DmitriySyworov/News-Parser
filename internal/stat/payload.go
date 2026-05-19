package stat

import (
	"time"
)

type ResponseStatCategoryDate struct {
	Categories []CategoryDbDate
}
type ResponseStatCategoryAll struct {
	Categories []CategoryDbAll
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
}
type ResponseStatArticleAll struct {
	Articles []ArticleDbAll
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
