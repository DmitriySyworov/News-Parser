package article_default

import "errors"

var (
	ErrLoadArticles            = errors.New("failed to load articles")
	ErrCategory                = errors.New("this category of articles does not exist")
	ErrIncorrectOnlyArticle    = errors.New("the 'isArticles' must be a boolean is true or false")
	ErrNotFoundArticle         = errors.New("article not found")
	ErrIncorrectUUIDArticle    = errors.New("article_uuid must be 36 characters long")
	ErrIncorrectIDArticleToday = errors.New("id must be 7 digits long")
)
