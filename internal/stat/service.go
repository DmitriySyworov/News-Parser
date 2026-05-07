package stat

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/pkg/event_bus"
	"log"
	"time"
)

type ServiceStat struct {
	*ServiceStatDep
	Repo *RepositoryStat
}
type ServiceStatDep struct {
	*event_bus.EventBus
}

func NewServiceStat(repo *RepositoryStat, dep *ServiceStatDep) *ServiceStat {
	return &ServiceStat{
		Repo:           repo,
		ServiceStatDep: dep,
	}
}
func (s *ServiceStat) GetStatCategoryByDate(dateStr string) (*ResponseStatCategoryDate, error) {
	date, errParse := time.Parse(time.DateOnly, dateStr)
	if errParse != nil {
		return nil, custom_errors.ErrIncorrectDate
	}
	place := 0
	if date == common.DateNow() {
		statCategoryToday, errGetToday := s.Repo.GetStatCategoryToday()
		if errGetToday != nil || len(statCategoryToday) == 0 {
			return nil, ErrStatLoad
		}
		var sliceRespDb []CategoryDbDate
		for _, stToday := range statCategoryToday {
			place++
			var dbResp CategoryDbDate
			category, ok := stToday.Member.(string)
			if !ok {
				return nil, ErrStatLoad
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
			return nil, errGetStat
		}
		return statCategories, nil
	}
}
func (s *ServiceStat) GetStatCategoryAllTime() (*ResponseStatCategoryAll, error) {
	allTimeStat, errGetAllTime := s.Repo.GetStatCategoryAllTime()
	if errGetAllTime != nil {
		return nil, ErrStatLoad
	}
	return allTimeStat, nil
}
func (s *ServiceStat) GetStatArticleByDate(dateStr string) (*ResponseStatArticleDate, error) {
	date, errParse := time.Parse(time.DateOnly, dateStr)
	if errParse != nil {
		return nil, custom_errors.ErrIncorrectDate
	}
	if date == common.DateNow() {
		statArticleToday, errGetToday := s.Repo.GetStatArticleToday()
		if errGetToday != nil || len(statArticleToday) == 0 {
			return nil, ErrStatLoad
		}
		var sliceRespDb []ArticleDbDate
		place := 0
		for _, stToday := range statArticleToday {
			place++
			var dbResp ArticleDbDate
			url, ok := stToday.Member.(string)
			if !ok {
				return nil, ErrStatLoad
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
			return nil, errGetStat
		}
		return statCategories, nil
	}
}
func (s *ServiceStat) GetStatArticleAllTime() (*ResponseStatArticleAll, error) {
	allTimeStat, errGetAllTime := s.Repo.GetStatArticleAllTime()
	if errGetAllTime != nil {
		return nil, ErrStatLoad
	}
	return allTimeStat, nil
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
					log.Println("failed to assertion type article, got: ", url)
				}
				s.Repo.addClickArticle(url)
			}
		}
	}
}
