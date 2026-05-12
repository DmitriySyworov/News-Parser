package article_user

import "errors"

var (
	ErrFailedToParse               = errors.New("failed to parse your link")
	ErrNotFoundRemoveArticles      = errors.New("no removed articles found")
	ErrIncorrectAllArticle         = errors.New("the 'allArticle' must be a boolean is true or false")
	ErrIncorrectAddText            = errors.New("the 'addText' must be a boolean is true or false")
	ErrIncorrectDeleteText         = errors.New("the 'deleteText' must be a boolean is true or false")
	ErrNotFoundUserArticle         = errors.New("no user article record found")
	ErrFailedRemoveArticle         = errors.New("failed to remove article entries")
	ErrIdAndAllArticleParams       = errors.New("the 'id' and 'allArticle' parameters cannot be specified simultaneously")
	ErrFailedRecoveryArticle       = errors.New("failed to recovery user article")
	ErrFailedUpdateUserArticle     = errors.New("failed to update user article")
	ErrDeleteTextAndAddTextTheSame = errors.New("'addText' and 'deleteText' parameters cannot be specified at the same time")
	ErrAllArticleAndDomain         = errors.New("if the 'allArticle' parameter is specified, the domain must be passed")
	ErrFailedParseText             = errors.New("failed to parse text")
	ErrIncorrectParams             = errors.New("the parameters are specified incorrectly. The action cannot be performed")
)
