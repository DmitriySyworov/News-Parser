package article

type ResponseCategoryToday struct {
	Header    string
	URL       string
	IDArticle uint
	Error     string
}
type ResponseCategoryArchive struct {
	Header      string
	URL         string
	UUIDArticle string
	Error       string
}
