package article_user

import "errors"

var (
	ErrFailedToParse          = errors.New("failed to parse your link")
	ErrNotFoundRemoveArticles = errors.New("no removed articles found")
	ErrIncorrectWithText      = errors.New("the 'withText' must be a boolean is true or false")
	ErrIncorrectAllArticle    = errors.New("the 'allArticle' must be a boolean is true or false")
	ErrIncorrectAddText       = errors.New("the 'addText' must be a boolean is true or false")
	ErrIncorrectArticleId     = errors.New("incorrect article_ID format or empty")
	ErrNotFoundUserArticle    = errors.New("no user article record found")
	ErrFailedRemoveArticle    = errors.New("failed to remove article entries")
	ErrIdAndAllArticleParams  = errors.New("the 'id' and 'allArticle' parameters cannot be specified simultaneously")
	ErrFailedRecoveryArticle  = errors.New("failed to recovery user article")
)
