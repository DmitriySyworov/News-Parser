package article_user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/model"
	"strconv"
	"sync"
	"time"
)

type ServiceArticleUser struct {
	Repo *RepositoryArticleUser
	Dep  *ServiceArticleUserDep
}
type ServiceArticleUserDep struct {
	di.IRepoUser
}

func NewServiceArticleUser(repo *RepositoryArticleUser, dep *ServiceArticleUserDep) *ServiceArticleUser {
	return &ServiceArticleUser{
		Repo: repo,
		Dep:  dep,
	}
}
func (s *ServiceArticleUser) CreateUserArticles(body *RequestCreateArticle, uuid, addTextStr string) (*ResponseUserArticles, error) {
	if _, errNotFound := s.Dep.IRepoUser.GetUserByUUID(uuid); errNotFound != nil {
		return nil, custom_errors.ErrUserNotExist
	}
	var isAddText bool
	if addTextStr == "false" || addTextStr == "" {
		isAddText = false
	} else if addTextStr == "true" {
		isAddText = true
	} else {
		return nil, ErrIncorrectAddText
	}
	var sliceUserArticle []model.UserArticle
	var wg sync.WaitGroup
	customParsing := NewCustomParsing(&wg, s.Repo)
	if isAddText {
		wg.Add(3)
		go customParsing.customParseCategory(body.URL, body.Category, uuid, isAddText)
		go customParsing.CustomParseArticle()
		go customParsing.createUserArticlesWithText()
		defer wg.Wait()
		for customLink := range customParsing.LinkUserCh {
			sliceUserArticle = append(sliceUserArticle, customLink)
		}
	} else {
		wg.Add(2)
		go customParsing.customParseCategory(body.URL, body.Category, uuid, isAddText)
		go customParsing.createUserArticlesWithoutText()
		defer wg.Wait()
		for customArticle := range customParsing.ArticleUserCh {
			sliceUserArticle = append(sliceUserArticle, customArticle)
		}
	}
	if len(sliceUserArticle) == 0 {
		return nil, ErrFailedToParse
	}
	return &ResponseUserArticles{SliceUserArticles: sliceUserArticle}, nil
}
func (s *ServiceArticleUser)UpdateUserArticle(body *RequestUpdateArticle, userUUID, idArticleStr, addTextStr, deleteTextStr, allArticleStr string)(*ResponseUserArticles, error){
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		return nil, custom_errors.ErrUserNotExist
	}
	if addTextStr != "" && deleteTextStr != ""{
		return nil,
	}
	if (allArticleStr != "" && body.Domain == "") || (allArticleStr == "" && body.Domain != ""){
		return nil,
	}
	var addText, deleteText, allArticle bool
	if addTextStr == "true" {
		addText = true
	} else if addTextStr == "false" || addTextStr == "" {
		addText = false
	} else {
		return nil,
	}
	if deleteTextStr == "true" {
		deleteText = true
	} else if deleteTextStr == "false" || deleteTextStr == "" {
		deleteText = false
	} else {
		return nil,
	}
	if allArticleStr == "true" {
		allArticle = true
	} else if allArticleStr == "false" ||allArticleStr  == "" {
		allArticle = false
	} else {
		return nil,
	}
	switch allArticle {
	case true:

	case false:
		idArticle, errParseId := strconv.Atoi(idArticleStr)
		if errParseId != nil {
			return nil, ErrIncorrectArticleId
		}
	userArticle, errGetUserArticle := s.Repo.GetUserArticle(userUUID, uint(idArticle))
	if errGetUserArticle != nil {
		return nil, ErrNotFoundUserArticle
	}
	if body.Category != "" &&  
	s.Repo.UpdateUserArticle(userUUID, &model.UserArticle{
		Header: ,
		URL: ,
		Text: ,
		Category: ,
		IDArticle: ,
		UUIDUser: ,
	})
	}
}
func (s *ServiceArticleUser) RemoveUserArticle(idArticleStr, userUUID, allArticleStr string) error {
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		return custom_errors.ErrUserNotExist
	}
	var isAllArticle bool
	if allArticleStr == "true" {
		isAllArticle = true
	} else if allArticleStr == "false" || allArticleStr == "" {
		isAllArticle = false
	} else {
		return ErrIncorrectAllArticle
	}
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		return ErrIncorrectArticleId
	}
	if isAllArticle {
		if errDeleteAll := s.Repo.DeleteAllUserArticle(userUUID); errDeleteAll != nil {
			return ErrFailedRemoveArticle
		}
	} else {
		if errDelete := s.Repo.RemoveUserArticleByID(userUUID, uint(idArticle)); errDelete != nil {
			return ErrFailedRemoveArticle
		}
	}
	return nil
}
func (s *ServiceArticleUser) GetUserArticle(userUUID, idArticleStr string) (*model.UserArticle, error) {
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		return nil, custom_errors.ErrUserNotExist
	}
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		return nil, ErrIncorrectArticleId
	}
	if userArticle, errGetUserArt := s.Repo.GetUserArticle(userUUID, uint(idArticle)); errGetUserArt != nil {
		return nil, ErrNotFoundUserArticle
	} else {
		return userArticle, nil
	}
}
func (s *ServiceArticleUser) GetAllUserArticles(userUUID, category, offsetStr, limitStr, withTextStr string) (*ResponseUserArticles, error) {
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		return nil, custom_errors.ErrUserNotExist
	}
	offset, limit, errValidateOffsetLimit := common.ValidateOffsetAndLimit(offsetStr, limitStr)
	if errValidateOffsetLimit != nil {
		return nil, errValidateOffsetLimit
	}
	var withText bool
	if withTextStr == "true" {
		withText = true
	} else if withTextStr == "false" || withTextStr == "" {
		withText = false
	} else {
		return nil, ErrIncorrectWithText
	}
	if withText {
		respUserArticles, errGetArticles := s.Repo.GetAllUserArticlesWithText(userUUID, category, offset, limit)
		if errGetArticles != nil {
			return nil, ErrNotFoundUserArticle
		}
		return respUserArticles, nil
	} else {
		respUserArticles, errGetArticles := s.Repo.GetAllUserArticlesWithoutText(userUUID, category, offset, limit)
		if errGetArticles != nil {
			return nil, ErrNotFoundUserArticle
		}
		return respUserArticles, nil
	}
}
func (s *ServiceArticleUser) GetRemoveUserArticle(userUUID, offsetStr, limitStr string) (*ResponseRemoveUserArticles, error) {
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		return nil, custom_errors.ErrUserNotExist
	}
	offset, limit, errValidateOffsetLimit := common.ValidateOffsetAndLimit(offsetStr, limitStr)
	if errValidateOffsetLimit != nil {
		return nil, errValidateOffsetLimit
	}
	sliceRemoveArticle, errGetRemoveArticles := s.Repo.GetRemoveUserArticle(userUUID, offset, limit)
	if errGetRemoveArticles != nil || len(sliceRemoveArticle) == 0 {
		return nil, ErrNotFoundRemoveArticles
	}
	for _, rmArticle := range sliceRemoveArticle {
		rmArticle.ExpiredAt = time.Unix(rmArticle.DeletedAt.Time.Unix()+common.UnixMonth, 0)
	}
	return &ResponseRemoveUserArticles{
		SliceRemoveUserArticles: sliceRemoveArticle,
	}, nil
}

func (s *ServiceArticleUser) RecoveryUserArticle(userUUID, idArticleStr, allArticleStr string) (*ResponseUserArticles, error) {
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		return nil, custom_errors.ErrUserNotExist
	}
	if allArticleStr != "" && idArticleStr != "" {
		return nil, ErrIdAndAllArticleParams
	}
	var isAllArticle bool
	if allArticleStr == "true" {
		isAllArticle = true
	} else if allArticleStr == "false" || allArticleStr == "" {
		isAllArticle = false
	} else {
		return nil, ErrIncorrectAllArticle
	}
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		return nil, ErrIncorrectArticleId
	}
	if !isAllArticle {
		userArticle, errRecoveryUserArticle := s.Repo.RecoveryUserArticle(userUUID, idArticle)
		if errRecoveryUserArticle != nil {
			return nil, ErrFailedRecoveryArticle
		}
		return &ResponseUserArticles{
			SliceUserArticles: []model.UserArticle{*userArticle},
		}, nil
	} else {
		sliceUserArticles, errRecoveryUserArticles := s.Repo.RecoveryAllUserArticle(userUUID)
		if errRecoveryUserArticles != nil {
			return nil, ErrFailedRecoveryArticle
		}
		return &ResponseUserArticles{
			SliceUserArticles: sliceUserArticles,
		}, nil
	}
}
func (s *ServiceArticleUser) DeletingRemoveUserArticle() {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	select {
	case <-ticker.C:
		s.Repo.deleteUserArticles()
	}
}
