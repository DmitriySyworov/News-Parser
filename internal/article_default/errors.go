package article_default

import "errors"

var (
	ErrLoadArticles         = errors.New("failed to load articles")
	ErrCategory             = errors.New("this category of articles does not exist")
	ErrChoiceArticlesFilter = errors.New("article_default filter incorrect or not specified")
)
