package article

import "errors"

var (
	ErrLoadArticles = errors.New("failed to load articles")
	ErrCategory     = errors.New("this category of articles does not exist")
)
