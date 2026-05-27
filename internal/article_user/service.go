package article_user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/model"
	"app/news-parser/internal/response"
	"app/news-parser/pkg/event_bus"
	"net/http"
	"sync"
	"time"
)

type ServiceArticleUser struct {
	Repo *RepositoryArticleUser
	Dep  *ServiceArticleUserDep
}
type ServiceArticleUserDep struct {
	di.IRepoUser
	*event_bus.EventBus
}

func NewServiceArticleUser(repo *RepositoryArticleUser, dep *ServiceArticleUserDep) *ServiceArticleUser {
	return &ServiceArticleUser{
		Repo: repo,
		Dep:  dep,
	}
}
func (s *ServiceArticleUser) CreateUserArticles(body *RequestCreateArticle, userUUID, addTextStr string) ([]ResponseCreateUserArticle, []response.Error) {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}
	var isAddText bool
	if addTextStr == "false" || addTextStr == "" {
		isAddText = false
	} else if addTextStr == "true" {
		isAddText = true
	} else {
		sliceError = append(sliceError, response.Error{
			Message: ErrIncorrectAddText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	var sliceUserArticle []ResponseCreateUserArticle
	customParsing := NewCustomParsing(s.Repo)
	counterArticle := 0
	if isAddText {
		go customParsing.customParseCategory(body.URL, body.Category, userUUID, isAddText)
		go customParsing.CustomParseArticle()
		go customParsing.createUserArticles(true)
		for customLink := range customParsing.RespUserCh {
			counterArticle++
			sliceUserArticle = append(sliceUserArticle, customLink)
		}
	} else {
		go customParsing.customParseCategory(body.URL, body.Category, userUUID, isAddText)
		go customParsing.createUserArticles(false)
		for customArticle := range customParsing.RespUserCh {
			counterArticle++
			sliceUserArticle = append(sliceUserArticle, customArticle)
		}
	}
	if len(sliceUserArticle) == 0 {
		sliceError = append(sliceError, response.Error{
			Message: ErrFailedToParse.Error(),
			Status:  http.StatusUnprocessableEntity,
		})
		return nil, sliceError
	}
	event := event_bus.Event{
		Name: common.EventCreateUserArticle,
		Data: common.StatDataUserArticle{
			UserUUID: userUUID,
			Number:   counterArticle,
		},
	}
	go s.Dep.EventBus.Publisher(&event)
	return sliceUserArticle, nil
}
func (s *ServiceArticleUser) UpdateUserArticle(category, userUUID, articleUUID, addTextStr, deleteTextStr string) (*model.UserArticle, []response.Error) {
	var sliceError []response.Error
	addText, deleteText, errSliceHelper := s.helperValidateUserAndAddTextAndDeleteText(userUUID, addTextStr, deleteTextStr)
	sliceError = append(sliceError, errSliceHelper...)
	userArticle, errGetUserArticle := s.Repo.GetUserArticle(userUUID, articleUUID)
	if errGetUserArticle != nil {
		sliceError = append(sliceError, response.Error{
			Message: ErrNotFoundUserArticle.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	}
	if len(userArticle.Text) <= 2 && deleteText {
		sliceError = append(sliceError, response.Error{
			Message: ErrDeleteText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(userArticle.Text) > 2 && addText {
		sliceError = append(sliceError, response.Error{
			Message: ErrAddText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	var text string
	var errParseText error
	if addText {
		text, errParseText = ParseText(userArticle.URL)
		if errParseText != nil || len(text) < 10 {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedParseText.Error(),
				Status:  http.StatusUnprocessableEntity,
			})
			return nil, sliceError
		}
	}
	event := event_bus.Event{
		Name: common.EventUpdateUserArticle,
		Data: common.StatDataUserArticle{
			UserUUID: userUUID,
			Number:   1,
		}}
	if category != "" && !addText && !deleteText {
		resUserArticle, errUpdateCategory := s.Repo.UpdateOneColumnUserArticle(userUUID, articleUUID, "category", category)
		if errUpdateCategory != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		go s.Dep.EventBus.Publisher(&event)
		return resUserArticle, nil
	} else if category == "" && !addText && deleteText {
		resUserArticle, errUpdateCategory := s.Repo.UpdateOneColumnUserArticle(userUUID, articleUUID, "text", "-")
		if errUpdateCategory != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		go s.Dep.EventBus.Publisher(&event)
		return resUserArticle, nil
	} else if category == "" && addText && !deleteText {
		resUserArticle, errUpdateCategory := s.Repo.UpdateOneColumnUserArticle(userUUID, articleUUID, "text", text)
		if errUpdateCategory != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		go s.Dep.EventBus.Publisher(&event)
		return resUserArticle, nil
	} else if category != "" && !addText && deleteText {
		resUserArticle := model.UserArticle{
			Header:      userArticle.Header,
			URL:         userArticle.URL,
			Text:        "-",
			Category:    category,
			ArticleUUID: userArticle.ArticleUUID,
			UserUUID:    userArticle.UserUUID,
		}
		errUpdateCategory := s.Repo.UpdateUserArticle(userUUID, resUserArticle.ArticleUUID, &resUserArticle)
		if errUpdateCategory != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		go s.Dep.EventBus.Publisher(&event)
		return &resUserArticle, nil
	} else if category != "" && addText && !deleteText {
		resUserArticle := model.UserArticle{
			Header:      userArticle.Header,
			URL:         userArticle.URL,
			Text:        text,
			Category:    category,
			ArticleUUID: userArticle.ArticleUUID,
			UserUUID:    userArticle.UserUUID,
		}
		errUpdateCategory := s.Repo.UpdateUserArticle(userUUID, resUserArticle.ArticleUUID, &resUserArticle)
		if errUpdateCategory != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedUpdateUserArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return nil, sliceError
		}
		go s.Dep.EventBus.Publisher(&event)
		return &resUserArticle, nil
	} else {
		return nil, []response.Error{
			{
				Message: ErrIncorrectParams.Error(),
				Status:  http.StatusBadRequest,
			}}
	}
}

func (s *ServiceArticleUser) UpdateBatchUserArticles(userUUID, domain, category, addTextStr, deleteTextStr string) ([]ResponseUserArticle, []response.Error) {
	var sliceError []response.Error
	addText, deleteText, errSliceHelper := s.helperValidateUserAndAddTextAndDeleteText(userUUID, addTextStr, deleteTextStr)
	sliceError = append(sliceError, errSliceHelper...)
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	var respUserArticles []ResponseUserArticle
	if category != "" && !addText && !deleteText {
		if !s.Repo.IsDomainArticleExist(userUUID, domain) {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		sliceUpdateArticles, countUpdate, errUpdateBatch := s.Repo.UpdateCategoryByDomainAll(userUUID, domain, category)
		if errUpdateBatch != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		for _, userArticle := range sliceUpdateArticles {
			respUserArticles = append(respUserArticles, ResponseUserArticle{
				Article: userArticle,
				SuccessOperation: SuccessOfTheOperation{
					Success: true,
					Message: "Update",
					Status:  http.StatusOK,
				},
			})
		}
		event := event_bus.Event{
			Name: common.EventUpdateUserArticle,
			Data: common.StatDataUserArticle{
				UserUUID: userUUID,
				Number:   int(countUpdate),
			},
		}
		go s.Dep.EventBus.Publisher(&event)
		return respUserArticles, nil
	} else if addText && !deleteText {
		sliceUserArticles, errGetUserArticlesByDomain := s.Repo.GetUserArticlesByDomain(userUUID, domain, false)
		if errGetUserArticlesByDomain != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		countUpdate := 0
		var wg sync.WaitGroup
		wg.Add(len(sliceUserArticles))
		for _, article := range sliceUserArticles {
			go func() {
				defer wg.Done()
				oldCategory := article.Category
				oldText := article.Text
				text, errParseText := ParseText(article.URL)
				if errParseText != nil || text == "" {
					respUserArticles = append(respUserArticles, ResponseUserArticle{
						Article: article,
						SuccessOperation: SuccessOfTheOperation{
							Success: false,
							Message: ErrFailedParseText.Error(),
							Status:  http.StatusUnprocessableEntity,
						},
					})
					return
				}
				if category != "" {
					article.Category = category
				}
				article.Text = text
				errUpdateArticle := s.Repo.UpdateUserArticle(userUUID, article.ArticleUUID, &article)
				if errUpdateArticle != nil {
					article.Category = oldCategory
					article.Text = oldText
					respUserArticles = append(respUserArticles, ResponseUserArticle{
						Article: article,
						SuccessOperation: SuccessOfTheOperation{
							Message: ErrFailedUpdateUserArticle.Error(),
							Status:  http.StatusInternalServerError,
						},
					})
					return
				}
				countUpdate++
				respUserArticles = append(respUserArticles, ResponseUserArticle{
					Article: article,
					SuccessOperation: SuccessOfTheOperation{
						Success: true,
						Message: "Update",
						Status:  http.StatusOK,
					},
				})
			}()
		}
		wg.Wait()
		event := event_bus.Event{
			Name: common.EventUpdateUserArticle,
			Data: common.StatDataUserArticle{
				UserUUID: userUUID,
				Number:   countUpdate,
			},
		}
		go s.Dep.EventBus.Publisher(&event)
		return respUserArticles, nil
	} else if !addText && deleteText {
		sliceUserArticles, errGetUserArticlesByDomain := s.Repo.GetUserArticlesByDomain(userUUID, domain, true)
		if errGetUserArticlesByDomain != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		countUpdate := 0
		for _, article := range sliceUserArticles {
			oldCategory := article.Category
			oldText := article.Text
			if category != "" {
				article.Category = category
			}
			article.Text = "-"
			errUpdateArticle := s.Repo.UpdateUserArticle(userUUID, article.ArticleUUID, &article)
			if errUpdateArticle != nil {
				article.Category = oldCategory
				article.Text = oldText
				respUserArticles = append(respUserArticles, ResponseUserArticle{
					Article: article,
					SuccessOperation: SuccessOfTheOperation{
						Success: false,
						Message: ErrFailedUpdateUserArticle.Error(),
						Status:  http.StatusInternalServerError,
					},
				})
			}
			countUpdate++
			respUserArticles = append(respUserArticles, ResponseUserArticle{
				Article: article,
				SuccessOperation: SuccessOfTheOperation{
					Success: true,
					Message: "Update",
					Status:  http.StatusOK,
				},
			})
		}
		event := event_bus.Event{
			Name: common.EventUpdateUserArticle,
			Data: common.StatDataUserArticle{
				UserUUID: userUUID,
				Number:   countUpdate,
			},
		}
		go s.Dep.EventBus.Publisher(&event)
		return respUserArticles, nil
	} else {
		return nil, []response.Error{
			{
				Message: ErrIncorrectParams.Error(),
				Status:  http.StatusBadRequest,
			}}
	}
}
func (s *ServiceArticleUser) helperValidateUserAndAddTextAndDeleteText(userUUID, addTextStr, deleteTextStr string) (bool, bool, []response.Error) {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}
	if addTextStr == "true" && deleteTextStr == "true" {
		sliceError = append(sliceError, response.Error{
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
		sliceError = append(sliceError, response.Error{
			Message: ErrIncorrectAddText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if deleteTextStr == "true" {
		deleteText = true
	} else if deleteTextStr == "false" || deleteTextStr == "" {
		deleteText = false
	} else {
		sliceError = append(sliceError, response.Error{
			Message: ErrIncorrectDeleteText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	return addText, deleteText, sliceError
}

const (
	typeSoftDelete = "soft-remove"
	typeHardDelete = "hard-remove"
)

func (s *ServiceArticleUser) RemoveUserArticle(articleUUID, userUUID, typeRemove string) []response.Error {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}
	if typeRemove == "" {
		typeRemove = typeSoftDelete
	}
	if typeRemove != typeSoftDelete && typeRemove != typeHardDelete {
		sliceError = append(sliceError, response.Error{
			Message: ErrTypeRemove.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if !s.Repo.IsUserArticleExistByUUID(userUUID, articleUUID) {
		sliceError = append(sliceError, response.Error{
			Message: ErrNotFoundUserArticle.Error(),
			Status:  http.StatusNotFound,
		})
	}
	if len(sliceError) != 0 {
		return sliceError
	}
	switch typeRemove {
	case typeSoftDelete:
		if !s.Repo.IsUserArticleExist(userUUID, articleUUID) {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return sliceError
		}
		if errRemove := s.Repo.RemoveUserArticleByUUID(userUUID, articleUUID); errRemove != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedRemoveArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return sliceError
		}
		event := event_bus.Event{
			Name: common.EventSoftDeleteUserArticle,
			Data: common.StatDataUserArticle{
				UserUUID: userUUID,
				Number:   1,
			},
		}
		go s.Dep.Publisher(&event)
	case typeHardDelete:
		if errDelete := s.Repo.DeleteUserArticleByUUID(userUUID, articleUUID); errDelete != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedRemoveArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return sliceError
		}
		event := event_bus.Event{
			Name: common.EventHardDeleteUserArticle,
			Data: common.StatDataUserArticle{
				UserUUID: userUUID,
				Number:   1,
			},
		}
		go s.Dep.Publisher(&event)
	}
	return nil
}

func (s *ServiceArticleUser) RemoveAllUserArticle(userUUID, typeRemove string) []response.Error {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}
	if typeRemove == "" {
		typeRemove = typeSoftDelete
	}
	if typeRemove != typeSoftDelete && typeRemove != typeHardDelete {
		sliceError = append(sliceError, response.Error{
			Message: ErrTypeRemove.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return sliceError
	}
	switch typeRemove {
	case typeSoftDelete:
		if !s.Repo.IsUserArticleExistNoRemoveAll(userUUID) {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return sliceError
		}
		countRemove, errRemoveAll := s.Repo.RemoveAllUserArticle(userUUID)
		if errRemoveAll != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedRemoveArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return sliceError
		}
		event := event_bus.Event{
			Name: common.EventSoftDeleteUserArticle,
			Data: common.StatDataUserArticle{
				UserUUID: userUUID,
				Number:   countRemove,
			},
		}
		go s.Dep.Publisher(&event)
	case typeHardDelete:
		if !s.Repo.IsUserArticleExistAll(userUUID) {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return sliceError
		}
		countDelete, errDeleteAll := s.Repo.DeleteAllUserArticle(userUUID)
		if errDeleteAll != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrFailedRemoveArticle.Error(),
				Status:  http.StatusInternalServerError,
			})
			return sliceError
		}
		event := event_bus.Event{
			Name: common.EventHardDeleteUserArticle,
			Data: common.StatDataUserArticle{
				UserUUID: userUUID,
				Number:   countDelete,
			},
		}
		go s.Dep.Publisher(&event)
	}
	return nil
}
func (s *ServiceArticleUser) GetUserArticle(userUUID, articleUUID string) (*model.UserArticle, []response.Error) {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}

	if len(sliceError) != 0 {
		return nil, sliceError
	}
	if userArticle, errGetUserArt := s.Repo.GetUserArticle(userUUID, articleUUID); errGetUserArt != nil {
		sliceError = append(sliceError, response.Error{
			Message: ErrNotFoundUserArticle.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	} else {
		return userArticle, nil
	}
}
func (s *ServiceArticleUser) GetAllUserArticles(userUUID, category, offsetStr, limitStr, withTextStr string) (*ResponseSliceUserArticles, []response.Error) {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
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
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrIncorrectWithText.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	switch withText {
	case true:
		respUserArticles, errGetArticles := s.Repo.GetAllUserArticlesWithText(userUUID, category, offset, limit)
		if errGetArticles != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		return respUserArticles, nil
	default:
		respUserArticles, errGetArticles := s.Repo.GetAllUserArticlesWithoutText(userUUID, category, offset, limit)
		if errGetArticles != nil {
			sliceError = append(sliceError, response.Error{
				Message: ErrNotFoundUserArticle.Error(),
				Status:  http.StatusNotFound,
			})
			return nil, sliceError
		}
		return respUserArticles, nil
	}
}
func (s *ServiceArticleUser) GetRemoveUserArticle(userUUID, offsetStr, limitStr string) (*ResponseRemoveUserArticles, []response.Error) {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
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
		sliceError = append(sliceError, response.Error{
			Message: ErrNotFoundRemoveArticles.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	}
	for i := range sliceRemoveArticle {
		sliceRemoveArticle[i].ExpiredAt = time.Unix(sliceRemoveArticle[i].DeletedAt.Time.Unix()+common.UnixMonth, 0)
	}
	return &ResponseRemoveUserArticles{
		SliceRemoveUserArticles: sliceRemoveArticle,
	}, nil
}

func (s *ServiceArticleUser) RecoveryUserArticle(userUUID, articleUUId string) (*model.UserArticle, []response.Error) {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}
	if !s.Repo.IsUserArticleRemoveExistByUUID(userUUID, articleUUId) {
		sliceError = append(sliceError, response.Error{
			Message: ErrNotFoundRemoveArticles.Error(),
			Status:  http.StatusNotFound,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	userArticle, errRecoveryUserArticle := s.Repo.RecoveryUserArticle(userUUID, articleUUId)
	if errRecoveryUserArticle != nil {
		sliceError = append(sliceError, response.Error{
			Message: ErrFailedRecoveryArticle.Error(),
			Status:  http.StatusInternalServerError,
		})
		return nil, sliceError
	}
	event := event_bus.Event{
		Name: common.EventRecoveryUserArticle,
		Data: common.StatDataUserArticle{
			UserUUID: userUUID,
			Number:   1,
		},
	}
	go s.Dep.Publisher(&event)
	return userArticle, nil
}
func (s *ServiceArticleUser) RecoveryAllUserArticle(userUUID string) (*ResponseSliceUserArticles, []response.Error) {
	var sliceError []response.Error
	if !s.Dep.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}
	if !s.Repo.IsUserArticleRecoveryExist(userUUID) {
		sliceError = append(sliceError, response.Error{
			Message: ErrNotFoundRemoveArticles.Error(),
			Status:  http.StatusNotFound,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	sliceUserArticles, countRecovery, errRecoveryUserArticles := s.Repo.RecoveryAllUserArticle(userUUID)
	if errRecoveryUserArticles != nil {
		sliceError = append(sliceError, response.Error{
			Message: ErrFailedRecoveryArticle.Error(),
			Status:  http.StatusInternalServerError,
		})
		return nil, sliceError
	}
	event := event_bus.Event{
		Name: common.EventRecoveryUserArticle,
		Data: common.StatDataUserArticle{
			UserUUID: userUUID,
			Number:   countRecovery,
		},
	}
	go s.Dep.Publisher(&event)
	return &ResponseSliceUserArticles{
		SliceUserArticles: sliceUserArticles,
	}, nil
}
func (s *ServiceArticleUser) DeletingRemoveUserArticle() {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	select {
	case <-ticker.C:
		s.Repo.deleteUserArticles()
	}
}
