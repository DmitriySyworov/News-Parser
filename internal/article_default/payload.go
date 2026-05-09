package article_default

type ResponseCategoryToday struct {
	Header    string
	URL       string
	Text      string
	IsArticle bool
	IDArticle uint
	Error     string
}
type ResponseCategoryArchive struct {
	Header      string
	URL         string
	UUIDArticle string
	Error       string
}
