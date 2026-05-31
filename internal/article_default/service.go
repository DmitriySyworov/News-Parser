package article_default

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/event_bus"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/model"
	"app/news-parser/internal/parsing_helper"
	"app/news-parser/internal/response"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ServiceArticle struct {
	repo *RepositoryArticle
	*ServiceArticleDep
}
type ServiceArticleDep struct {
	*event_bus.EventBus
	ResBrowser *parsing_helper.Browser
	*loggers.Logger
	di.IRepoStat
}

const (
	food        = "food"
	politics    = "politics"
	sport       = "sport"
	cloth       = "cloth"
	business    = "business"
	electronics = "electronics"

	lengthIdArticle = 7
)

var StorageCategories = []string{food, politics, sport, cloth, business, electronics}

func NewServiceArticle(repoArticle *RepositoryArticle, dep *ServiceArticleDep) *ServiceArticle {
	return &ServiceArticle{
		repo:              repoArticle,
		ServiceArticleDep: dep,
	}
}
func (s *ServiceArticle) GetArticlesInCategoryToday(category, offsetStr, limitStr, filterStr, withTextStr string) ([]ResponseCategoryToday, []response.Error) {
	var sliceError []response.Error
	if !validateCategories(category) {
		sliceError = append(sliceError, response.Error{
			Message: ErrCategory.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	offset, limit, errValidateOffsetLimit := common.ValidateOffsetAndLimit(offsetStr, limitStr)
	if errValidateOffsetLimit != nil {
		sliceError = append(sliceError, errValidateOffsetLimit...)
	}
	var filter, withText bool
	if filterStr == "true" || filterStr == "" {
		filter = true
	} else if filterStr == "false" {
		filter = false
	} else {
		sliceError = append(sliceError, response.Error{
			Message: ErrIncorrectOnlyArticle.Error(),
			Status:  http.StatusBadRequest,
		})

	}
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
	allArticle, errGetAllArticle := s.repo.GetArticlesInCategoryToday(category, offset, limit, filter, withText)
	if errGetAllArticle != nil {
		sliceError = append(sliceError, response.Error{
			Message: ErrNotFoundArticle.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	}
	go s.Publisher(&event_bus.Event{
		Name: event_bus.EventClickCategory,
		Data: category,
	})
	return allArticle, nil
}
func (s *ServiceArticle) GetArticleToday(idStr string) (*model.ArticleToday, *response.Error) {
	id, errParseId := strconv.Atoi(idStr)
	if errParseId != nil {
		return nil, &response.Error{
			Message: ErrIncorrectIDArticleToday.Error(),
			Status:  http.StatusBadRequest,
		}
	}
	article, errGetArticle := s.repo.GetArticleToday(id)
	if errGetArticle != nil {
		return nil, &response.Error{
			Message: ErrNotFoundArticle.Error(),
			Status:  http.StatusNotFound,
		}
	}
	go s.Publisher(&event_bus.Event{
		Name: event_bus.EventClickArticle,
		Data: article.URL,
	})
	return article, nil

}
func (s *ServiceArticle) GetArticlesInCategoryArchive(category, offsetStr, limitStr, dateStr string) ([]ResponseCategoryArchive, []response.Error) {
	var sliceError []response.Error
	if !validateCategories(category) {
		sliceError = append(sliceError, response.Error{
			Message: ErrCategory.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	offset, limit, errOffsetLimit := common.ValidateOffsetAndLimit(offsetStr, limitStr)
	if len(errOffsetLimit) != 0 {
		sliceError = append(sliceError, errOffsetLimit...)
	}
	date, errParseDate := time.Parse(time.DateOnly, dateStr)
	if errParseDate != nil {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrIncorrectDate.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	archiveArticles, errGetArticlesArch := s.repo.GetArticlesInCategoryArchive(category, offset, limit, date)
	if errGetArticlesArch != nil {
		sliceError = append(sliceError, response.Error{
			Message: errGetArticlesArch.Error(),
			Status:  http.StatusNotFound,
		})
		return nil, sliceError
	}
	var respCategoryArch []ResponseCategoryArchive
	for _, arch := range archiveArticles {
		var tempArch ResponseCategoryArchive
		tempArch.UUIDArticle = arch.ArticleArchiveUUID
		tempArch.URL = arch.URL
		tempArch.Header = arch.Header
		respCategoryArch = append(respCategoryArch, tempArch)
	}
	go s.Publisher(&event_bus.Event{
		Name: event_bus.EventClickCategory,
		Data: category,
	})
	return respCategoryArch, nil
}
func (s *ServiceArticle) GetArchiveArticle(uuid string) (*model.ArticleArchive, *response.Error) {
	archArticle, errGetArchArticle := s.repo.GetArchiveArticle(uuid)
	if errGetArchArticle != nil {
		return nil, &response.Error{
			Message: ErrNotFoundArticle.Error(),
			Status:  http.StatusNotFound,
		}
	}
	go s.Publisher(&event_bus.Event{
		Name: event_bus.EventClickArticle,
		Data: archArticle.URL,
	})
	return archArticle, nil
}

func validateCategories(category string) bool {
	for _, c := range StorageCategories {
		if category == c {
			return true
		}
	}
	return false
}
func (s *ServiceArticle) ReplacementInfo() {
	tickerEveryDay := time.NewTicker(common.Day)
	tickerCheckEmpty := time.NewTicker(20 * time.Second) //!!!
	defer tickerEveryDay.Stop()
	defer tickerCheckEmpty.Stop()
	for {
		select {
		case <-tickerEveryDay.C:
			go s.saveStatCategoryRedis()
			go s.saveStatArticleRedis()
			go s.saveDataRedis()
			go s.loadNewInfo()
		case <-tickerCheckEmpty.C:
			if !s.repo.isRedisArticleExist() {
				go s.loadNewInfo()
			}
		}
	}
}
func (s *ServiceArticle) saveStatCategoryRedis() {
	statCategory, errGetStatC := s.IRepoStat.GetStatCategoryToday()
	if errGetStatC != nil {
		log.Println(errGetStatC)
	}
	for _, statC := range statCategory {
		category, ok := statC.Member.(string)
		if !ok {
			log.Println("failed to assertion category, got: ", statC.Member)
			continue
		}
		errCreateStatCategory := s.IRepoStat.CreateStatCategory(&model.CategoryStat{
			Category: category,
			Click:    uint(statC.Score),
			Date:     common.DateNow(),
		})
		if errCreateStatCategory != nil {
			log.Println(errCreateStatCategory)
		}
	}
}
func (s *ServiceArticle) saveStatArticleRedis() {
	statArticle, errGetStatA := s.IRepoStat.GetStatArticleToday()
	if errGetStatA != nil {
		log.Println(errGetStatA)
	}
	for _, statA := range statArticle {
		url, ok := statA.Member.(string)
		if !ok {
			log.Println("failed to assertion category, got: ", statA.Member)
			continue
		}
		errCreateStatCategory := s.IRepoStat.CreateStatArticle(&model.ArticleStat{
			URL:   url,
			Click: uint(statA.Score),
			Date:  common.DateNow(),
		})
		if errCreateStatCategory != nil {
			log.Println(errCreateStatCategory)
		}
	}
}

func (s *ServiceArticle) saveDataRedis() {
	sliceArticles, errAllArticle := s.repo.allArticlesRedis()
	if errAllArticle != nil {
		log.Println(errAllArticle)
	}
	for _, article := range sliceArticles {
		if errCreate := s.repo.CreateArchiveArticle(&article); errCreate != nil {
			log.Println(errCreate)
		}
	}
}

func (s *ServiceArticle) loadNewInfo() {
	var linksList = [3]string{
		`https://www.bbc.com/russian>politics>https://www.bbc.com>Время чтения:>Темы>articles>~`,
		`https://edition.cnn.com/sport>sport>https://edition.cnn.com>By>World>/\d\d\d\d/\d\d/\d\d/sport>Video`,
		`https://www.cnbc.com/business/>business>https://www.cnbc.com>In this article>Choose CNBC as your preferred source on Google and never miss a moment from the most trusted name in business news.>/\d\d\d\d/\d\d/\d\d/>~`,
	}
	var wg sync.WaitGroup
	for _, list := range linksList {
		data := strings.Split(list, ">")
		parse := NewParsing(&wg, s.repo, s.ServiceArticleDep.ResBrowser, s.Logger)
		wg.Add(3)
		go func(url, category, domain, flagText, isArticleOnHeader string) {
			defer wg.Done()
			parse.parseCategory(url, category, domain, flagText, isArticleOnHeader)
			close(parse.LinkCh)
		}(data[0], data[1], data[2], data[5], data[6])

		go func(domain, startWord, stopWord string) {
			defer wg.Done()
			parse.parseArticle(domain, startWord, stopWord)
			close(parse.IsOk)
			close(parse.ArticleCh)
		}(data[2], data[3], data[4])

		go func(category string) {
			defer wg.Done()
			parse.createRdb(category)
		}(data[1])
	}
	wg.Wait()
}
func (s *ServiceArticle) RemoveUserArticles() {
	ticker := time.NewTicker(common.Day)
	defer ticker.Stop()
	select {
	case <-ticker.C:

	}
}
