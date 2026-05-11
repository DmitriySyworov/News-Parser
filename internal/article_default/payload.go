package article_default

type ResponseCategoryToday struct {
	Header    string
	URL       string
	Text      string
	IsArticle bool
	IDArticle uint
}
type ResponseCategoryArchive struct {
	Header      string
	URL         string
	UUIDArticle string
}
