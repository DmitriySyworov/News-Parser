package article_user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/model"
	"fmt"
	"net/http"
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
func (s *ServiceArticleUser) CreateUserArticles(body *RequestCreateArticle, uuid, addTextStr string) (*ResponseSliceUserArticles, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(uuid) {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusUnauthorized,
		})
	}
	var isAddText bool
	if addTextStr == "false" || addTextStr == "" {
		isAddText = false
	} else if addTextStr == "true" {
		isAddText = true
	} else {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrIncorrectAddText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
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
			fmt.Println(customArticle)
			sliceUserArticle = append(sliceUserArticle, customArticle)
		}
	}
	if len(sliceUserArticle) == 0 {
		sliceError = append(sliceError, custom_errors.Error{Message: ErrFailedToParse.Error(), Status: http.StatusUnprocessableEntity})
		return nil, sliceError
	}
	return &ResponseSliceUserArticles{SliceUserArticles: sliceUserArticle}, nil
}
func (s *ServiceArticleUser) UpdateUserArticle(category, userUUID, idArticleStr, addTextStr, deleteTextStr string) (*model.UserArticle, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	addText, deleteText, errSliceHelper := s.helperValidateUserAndAddTextAndDeleteText(userUUID, addTextStr, deleteTextStr)
	sliceError = append(sliceError, errSliceHelper...)
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrIncorrectArticleId.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	userArticle, errGetUserArticle := s.Repo.GetUserArticle(userUUID, uint(idArticle))
	if errGetUserArticle != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrNotFoundUserArticle.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	}
	if category != "" && !addText && deleteText {
		text, errParseText := ParseText(userArticle.URL)
		if errParseText != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedParseText.Error(),
				Status:  http.StatusUnprocessableEntity,
			})
			return nil, sliceError
		}
		resUserArticle := model.UserArticle{
			Header:    userArticle.Header,
			URL:       userArticle.URL,
			Text:      text,
			Category:  category,
			IDArticle: userArticle.IDArticle,
			UUIDUser:  userArticle.UUIDUser,
		}
		errUpdateCategory := s.Repo.UpdateUserArticle(userUUID, &resUserArticle)
		if errUpdateCategory != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		return &resUserArticle, nil
	} else if category != "" && addText && !deleteText {
		resUserArticle := model.UserArticle{
			Header:    userArticle.Header,
			URL:       userArticle.URL,
			Text:      "-",
			Category:  category,
			IDArticle: userArticle.IDArticle,
			UUIDUser:  userArticle.UUIDUser,
		}
		errUpdateCategory := s.Repo.UpdateUserArticle(userUUID, &resUserArticle)
		if errUpdateCategory != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		return &resUserArticle, nil
	} else if category != "" && !addText && !deleteText {
		resUserArticle, errUpdateCategory := s.Repo.UpdateOneColumnUserArticle(userUUID, category, "category")
		if errUpdateCategory != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		return resUserArticle, nil
	} else if category == "" && !addText && deleteText {
		resUserArticle, errUpdateCategory := s.Repo.UpdateOneColumnUserArticle(userUUID, "-", "text")
		if errUpdateCategory != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		return resUserArticle, nil
	} else if category == "" && addText && !deleteText {
		text, errParseText := ParseText(userArticle.URL)
		if errParseText != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedParseText.Error(),
				Status:  http.StatusUnprocessableEntity,
			})
			return nil, sliceError
		}
		resUserArticle, errUpdateCategory := s.Repo.UpdateOneColumnUserArticle(userUUID, text, "text")
		if errUpdateCategory != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		return resUserArticle, nil
	} else {
		return nil, []custom_errors.Error{
			{
				Message: ErrIncorrectParams.Error(),
				Status:  http.StatusBadRequest,
			}}
	}
}
func (s *ServiceArticleUser) UpdateBatchUserArticles(domain, userUUID, addTextStr, deleteTextStr string) ([]ResponseUserArticle, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	addText, deleteText, errSliceHelper := s.helperValidateUserAndAddTextAndDeleteText(userUUID, addTextStr, deleteTextStr)
	sliceError = append(sliceError, errSliceHelper...)
	if len(sliceError) != 0 {
		return nil, sliceError
	}

	var respUserArticles []ResponseUserArticle
	if addText && !deleteText {
		sliceUserArticles, errGetUserArticlesByDomain := s.Repo.GetUserArticlesByDomain(userUUID, domain, true)
		if errGetUserArticlesByDomain != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		for _, article := range sliceUserArticles {
			text, errParseText := ParseText(article.URL)
			if errParseText != nil {
				respUserArticles = append(respUserArticles, ResponseUserArticle{
					Article: article,
					Error:   ErrFailedParseText.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
			}
			updateArticle, errUpdateArticle := s.Repo.UpdateOneColumnUserArticle(userUUID, text, "text")
			if errUpdateArticle != nil {
				respUserArticles = append(respUserArticles, ResponseUserArticle{
					Article: article,
					Error:   ErrFailedUpdateUserArticle.Error(),
					Status:  http.StatusInternalServerError,
				})
			} else {
				respUserArticles = append(respUserArticles, ResponseUserArticle{
					Article: *updateArticle,
					Status:  http.StatusOK,
				})
			}
		}
	} else if !addText && deleteText {
		sliceUserArticles, errGetUserArticlesByDomain := s.Repo.GetUserArticlesByDomain(userUUID, domain, true)
		if errGetUserArticlesByDomain != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		for _, article := range sliceUserArticles {
			updateArticle, errUpdateArticle := s.Repo.UpdateOneColumnUserArticle(userUUID, "-", "text")
			if errUpdateArticle != nil {
				respUserArticles = append(respUserArticles, ResponseUserArticle{
					Article: article,
					Error:   ErrFailedUpdateUserArticle.Error(),
					Status:  http.StatusInternalServerError,
				})
			} else {
				respUserArticles = append(respUserArticles, ResponseUserArticle{
					Article: *updateArticle,
					Status:  http.StatusOK,
				})
			}
		}
	}
	return respUserArticles, nil
}
func (s *ServiceArticleUser) helperValidateUserAndAddTextAndDeleteText(userUUID, addTextStr, deleteTextStr string) (bool, bool, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusUnauthorized,
		})
	}
	if addTextStr != "" && deleteTextStr != "" {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrDeleteTextAndAddTextTheSame.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	var addText, deleteText bool
	if addTextStr == "true" {
		addText = true
	} else if addTextStr == "false" || addTextStr == "" {
		addText = false
	} else {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrIncorrectAddText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if deleteTextStr == "true" {
		deleteText = true
	} else if deleteTextStr == "false" || deleteTextStr == "" {
		deleteText = false
	} else {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrIncorrectDeleteText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	return addText, deleteText, sliceError
}
func (s *ServiceArticleUser) RemoveUserArticle(idArticleStr, userUUID, allArticleStr string) []custom_errors.Error {
	var sliceError []custom_errors.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusUnauthorized,
		})
	}
	var isAllArticle bool
	if allArticleStr == "true" {
		isAllArticle = true
	} else if allArticleStr == "false" || allArticleStr == "" {
		isAllArticle = false
	} else {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrIncorrectAllArticle.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrIncorrectArticleId.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return sliceError
	}
	if isAllArticle {
		if errDeleteAll := s.Repo.DeleteAllUserArticle(userUUID); errDeleteAll != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedRemoveArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return sliceError
		}
	} else {
		if errDelete := s.Repo.RemoveUserArticleByID(userUUID, uint(idArticle)); errDelete != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedRemoveArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return sliceError
		}
	}
	return nil
}
func (s *ServiceArticleUser) GetUserArticle(userUUID, idArticleStr string) (*model.UserArticle, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrIncorrectArticleId.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	if userArticle, errGetUserArt := s.Repo.GetUserArticle(userUUID, uint(idArticle)); errGetUserArt != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrNotFoundUserArticle.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	} else {
		return userArticle, nil
	}
}
func (s *ServiceArticleUser) GetAllUserArticles(userUUID, category, offsetStr, limitStr, withTextStr string) (*ResponseSliceUserArticles, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusUnauthorized,
		})
	}
	offset, limit, errValidateOffsetLimit := common.ValidateOffsetAndLimit(offsetStr, limitStr)
	sliceError = append(sliceError, errValidateOffsetLimit...)
	var withText bool
	if withTextStr == "true" {
		withText = true
	} else if withTextStr == "false" || withTextStr == "" {
		withText = false
	} else {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrIncorrectWithText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	if withText {
		respUserArticles, errGetArticles := s.Repo.GetAllUserArticlesWithText(userUUID, category, offset, limit)
		if errGetArticles != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusBadRequest,
			})
			return nil, sliceError
		}
		return respUserArticles, nil
	} else {
		respUserArticles, errGetArticles := s.Repo.GetAllUserArticlesWithoutText(userUUID, category, offset, limit)
		if errGetArticles != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusBadRequest,
			})
			return nil, sliceError
		}
		return respUserArticles, nil
	}
}
func (s *ServiceArticleUser) GetRemoveUserArticle(userUUID, offsetStr, limitStr string) (*ResponseRemoveUserArticles, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusUnauthorized,
		})
	}
	offset, limit, errValidateOffsetLimit := common.ValidateOffsetAndLimit(offsetStr, limitStr)
	if errValidateOffsetLimit != nil {
		sliceError = append(sliceError, errValidateOffsetLimit...)
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	sliceRemoveArticle, errGetRemoveArticles := s.Repo.GetRemoveUserArticle(userUUID, offset, limit)
	if errGetRemoveArticles != nil || len(sliceRemoveArticle) == 0 {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrNotFoundRemoveArticles.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	}
	for _, rmArticle := range sliceRemoveArticle {
		rmArticle.ExpiredAt = time.Unix(rmArticle.DeletedAt.Time.Unix()+common.UnixMonth, 0)
	}
	return &ResponseRemoveUserArticles{
		SliceRemoveUserArticles: sliceRemoveArticle,
	}, nil
}

func (s *ServiceArticleUser) RecoveryUserArticle(userUUID, idArticleStr, allArticleStr string) (*ResponseSliceUserArticles, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusUnauthorized,
		})
	}
	if allArticleStr != "" && idArticleStr != "" {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrIdAndAllArticleParams.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	var isAllArticle bool
	if allArticleStr == "true" {
		isAllArticle = true
	} else if allArticleStr == "false" || allArticleStr == "" {
		isAllArticle = false
	} else {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrIncorrectAllArticle.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrIncorrectArticleId.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	if !isAllArticle {
		userArticle, errRecoveryUserArticle := s.Repo.RecoveryUserArticle(userUUID, idArticle)
		if errRecoveryUserArticle != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedRecoveryArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		return &ResponseSliceUserArticles{
			SliceUserArticles: []model.UserArticle{*userArticle},
		}, nil
	} else {
		sliceUserArticles, errRecoveryUserArticles := s.Repo.RecoveryAllUserArticle(userUUID)
		if errRecoveryUserArticles != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: ErrFailedRecoveryArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		return &ResponseSliceUserArticles{
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
