package stat

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/pkg/event_bus"
	"log"
	"net/http"
	"time"
)

type ServiceStat struct {
	*ServiceStatDep
	Repo *RepositoryStat
}
type ServiceStatDep struct {
	*event_bus.EventBus
	di.IRepoUser
}

func NewServiceStat(repo *RepositoryStat, dep *ServiceStatDep) *ServiceStat {
	return &ServiceStat{
		Repo:           repo,
		ServiceStatDep: dep,
	}
}
func (s *ServiceStat) GetStatCategoryByDate(dateStr string) (*ResponseStatCategoryDate, *custom_errors.Error) {
	date, errParse := time.Parse(time.DateOnly, dateStr)
	if errParse != nil {
		return nil, &custom_errors.Error{
			Message: ErrIncorrectDate.Error(),
			Status:  http.StatusBadRequest,
		}
	}
	place := 0
	if date == common.DateNow() {
		statCategoryToday, errGetToday := s.Repo.GetStatCategoryToday()
		if errGetToday != nil || len(statCategoryToday) == 0 {
			return nil, &custom_errors.Error{
				Message: ErrStatNotFound.Error(),
				Status:  http.StatusNotFound,
			}
		}
		var sliceRespDb []CategoryDbDate
		for _, stToday := range statCategoryToday {
			place++
			var dbResp CategoryDbDate
			category, ok := stToday.Member.(string)
			if !ok {
				return nil, &custom_errors.Error{
					Message: ErrStatNotFound.Error(),
					Status:  http.StatusNotFound,
				}
			}
			dbResp.Category = category
			dbResp.Date = date
			dbResp.Click = uint(stToday.Score)
			dbResp.Place = uint(place)
			sliceRespDb = append(sliceRespDb, dbResp)
		}
		return &ResponseStatCategoryDate{
			Categories: sliceRespDb,
		}, nil
	} else {
		statCategories, errGetStat := s.Repo.GetStatCategoryByDate(date)
		if errGetStat != nil {
			return nil, &custom_errors.Error{
				Message: ErrStatNotFound.Error(),
				Status:  http.StatusNotFound,
			}
		}
		return statCategories, nil
	}
}
func (s *ServiceStat) GetStatArticleByDate(dateStr string) (*ResponseStatArticleDate, *custom_errors.Error) {
	date, errParse := time.Parse(time.DateOnly, dateStr)
	if errParse != nil {
		return nil, &custom_errors.Error{
			Message: ErrIncorrectDate.Error(),
			Status:  http.StatusBadRequest,
		}
	}
	if date == common.DateNow() {
		statArticleToday, errGetToday := s.Repo.GetStatArticleToday()
		if errGetToday != nil || len(statArticleToday) == 0 {
			return nil, &custom_errors.Error{
				Message: ErrStatNotFound.Error(),
				Status:  http.StatusNotFound,
			}
		}
		var sliceRespDb []ArticleDbDate
		place := 0
		for _, stToday := range statArticleToday {
			place++
			var dbResp ArticleDbDate
			url, ok := stToday.Member.(string)
			if !ok {
				return nil, &custom_errors.Error{
					Message: ErrStatNotFound.Error(),
					Status:  http.StatusNotFound,
				}
			}
			dbResp.URL = url
			dbResp.Date = date
			dbResp.Click = uint(stToday.Score)
			dbResp.Place = uint(place)
			sliceRespDb = append(sliceRespDb, dbResp)
		}
		return &ResponseStatArticleDate{
			Articles: sliceRespDb,
		}, nil
	} else {
		statCategories, errGetStat := s.Repo.GetStatArticleByDate(date)
		if errGetStat != nil {
			return nil, &custom_errors.Error{
				Message: ErrStatNotFound.Error(),
				Status:  http.StatusNotFound,
			}
		}
		return statCategories, nil
	}
}
func (s *ServiceStat) GetUserArticleStat(userUUID, dateStr string) (*ResponseUserArticleStat, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !s.IRepoUser.IsUserExistByUUID(userUUID) {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		})
	}
	date, errParse := time.Parse(time.DateOnly, dateStr)
	if errParse != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrIncorrectDate.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	statUser, errGetStat := s.Repo.GetUserArticleStat(userUUID, date)
	if errGetStat != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrStatNotFound.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	}
	return statUser, nil
}
func (s *ServiceStat) GetUserArticleAllTimeStat(userUUID string) (*ResponseUserArticleAllTimeStat, *custom_errors.Error) {
	if !s.IRepoUser.IsUserExistByUUID(userUUID) {
		return nil, &custom_errors.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		}
	}
	respAllTimeStat, errGetStat := s.Repo.GetUserArticleAllTimeStat(userUUID)
	if errGetStat != nil {
		return nil, &custom_errors.Error{
			Message: ErrStatNotFound.Error(),
			Status:  http.StatusNotFound,
		}
	}
	return respAllTimeStat, nil
}
func (s *ServiceStat) PushInStat() {
	for {
		for event := range s.Subscriber() {
			switch event.Name {
			case common.EventClickCategory:
				category, ok := event.Data.(string)
				if !ok {
					log.Println("failed to assertion type category, got: ", category)
				}
				s.Repo.addClickCategory(category)
			case common.EventClickArticle:
				url, ok := event.Data.(string)
				if !ok {
					log.Println("failed to assertion type article_default, got: ", url)
				}
				s.Repo.addClickArticle(url)
			case common.EventCreateUserArticle:
				statData, ok := event.Data.(common.StatDataUserArticle)
				if !ok {
					log.Println("failed to assertion type statUserData, got: ")
				}
				s.Repo.addCreateUserArticle(statData.Number, statData.UserUUID)
			case common.EventUpdateUserArticle:
				statData, ok := event.Data.(common.StatDataUserArticle)
				if !ok {
					log.Println("failed to assertion type statUserData, got: ")
				}
				s.Repo.addUpdateUserArticle(statData.Number, statData.UserUUID)
			case common.EventSoftDeleteUserArticle:
				statData, ok := event.Data.(common.StatDataUserArticle)
				if !ok {
					log.Println("failed to assertion type statUserData, got: ")
				}
				s.Repo.addSoftDeleteUserArticle(statData.Number, statData.UserUUID)
			case common.EventHardDeleteUserArticle:
				statData, ok := event.Data.(common.StatDataUserArticle)
				if !ok {
					log.Println("failed to assertion type statUserData, got: ")
				}
				s.Repo.addHardDeleteUserArticle(statData.Number, statData.UserUUID)
			case common.EventRecoveryUserArticle:
				statData, ok := event.Data.(common.StatDataUserArticle)
				if !ok {
					log.Println("failed to assertion type statUserData, got: ")
				}
				s.Repo.addRecoveryUserArticle(statData.Number, statData.UserUUID)
			}
		}
	}
}
