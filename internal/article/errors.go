package article

import "errors"

var (
	ErrLoadArticles         = errors.New("failed to load articles")
	ErrCategory             = errors.New("this category of articles does not exist")
	ErrIncorrectId          = errors.New("incorrect article ID format or empty")
	ErrIncorrectLimit       = errors.New("incorrect format limit or empty")
	ErrChoiceArticlesFilter = errors.New("article filter incorrect or not specified")
	ErrFailedToParse        = errors.New("failed to parse your link")
)
