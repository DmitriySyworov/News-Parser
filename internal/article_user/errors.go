package article_user

import "errors"

var (
	ErrParseInitialization         = errors.New("parser initialization error")
	ErrDeleteText                  = errors.New("no text anymore")
	ErrAddText                     = errors.New("text has already been added")
	ErrFailedToParse               = errors.New("failed to parse your link")
	ErrNotFoundRemoveArticles      = errors.New("no removed article found")
	ErrIncorrectAddText            = errors.New("the 'addText' must be a boolean is true or false")
	ErrIncorrectDeleteText         = errors.New("the 'deleteText' must be a boolean is true or false")
	ErrNotFoundUserArticle         = errors.New("no user article record found")
	ErrFailedRemoveArticle         = errors.New("failed to remove article entries")
	ErrFailedRecoveryArticle       = errors.New("failed to recovery user article")
	ErrFailedUpdateUserArticle     = errors.New("failed to update user article")
	ErrDeleteTextAndAddTextTheSame = errors.New("'addText' and 'deleteText' parameters cannot be specified at the same time")
	ErrAllArticleAndDomain         = errors.New("if the 'allArticle' parameter is specified, the domain must be passed")
	ErrFailedParseText             = errors.New("failed to parse text")
	ErrIncorrectParams             = errors.New("the parameters are specified incorrectly. The action cannot be performed")
	ErrUserURL                     = errors.New("could not get a response from the resource you passed")
	ErrFailedParseBody             = errors.New("failed to get the body of the main page of the resource")
	ErrSaveParseData               = errors.New("failed to save data received data")
	ErrTypeRemove                  = errors.New("the 'type' must be a soft-remove or hard-remove")
)
