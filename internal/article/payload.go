package article

type ResponseCategoryToday struct {
	Header    string
	URL       string
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
