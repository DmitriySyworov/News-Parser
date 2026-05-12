package article_default

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/model"
	"app/news-parser/pkg/event_bus"
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
func (s *ServiceArticle) GetArticlesInCategoryToday(category, offsetStr, limitStr, filterStr, withTextStr string) ([]ResponseCategoryToday, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !validateCategories(category) {
		sliceError = append(sliceError, custom_errors.Error{
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
		sliceError = append(sliceError, custom_errors.Error{
			Message: ErrIncorrectOnlyArticle.Error(),
			Status:  http.StatusBadRequest,
		})

	}
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
	allArticle, errGetAllArticle := s.repo.GetArticlesInCategoryToday(category, offset, limit, filter, withText)
	if errGetAllArticle != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: errGetAllArticle.Error(),
			Status:  http.StatusInternalServerError,
		})
		return nil, sliceError
	}
	go s.Publisher(&event_bus.Event{
		Name: common.EventClickCategory,
		Data: category,
	})
	return allArticle, nil
}
func (s *ServiceArticle) GetArticleToday(idStr string) (*model.ArticleToday, *custom_errors.Error) {
	id, errParseId := strconv.Atoi(idStr)
	if errParseId != nil {
		return nil, &custom_errors.Error{
			Message: custom_errors.ErrIncorrectArticleId.Error(),
			Status:  http.StatusBadRequest,
		}
	}
	article, errGetArticle := s.repo.GetArticleToday(id)
	if errGetArticle != nil {
		return nil, &custom_errors.Error{
			Message: ErrLoadArticles.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	go s.Publisher(&event_bus.Event{
		Name: common.EventClickArticle,
		Data: article.URL,
	})
	return article, nil

}
func (s *ServiceArticle) GetArticlesInCategoryArchive(category, offsetStr, limitStr, dateStr string) ([]ResponseCategoryArchive, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	if !validateCategories(category) {
		sliceError = append(sliceError, custom_errors.Error{
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
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrIncorrectDate.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if len(sliceError) != 0 {
		return nil, sliceError
	}
	archiveArticles, errGetArticlesArch := s.repo.GetArticlesInCategoryArchive(category, offset, limit, date)
	if errGetArticlesArch != nil {
		sliceError = append(sliceError, custom_errors.Error{
			Message: errGetArticlesArch.Error(),
			Status:  http.StatusUnauthorized,
		})
		return nil, sliceError
	}
	var respCategoryArch []ResponseCategoryArchive //!!!!!
	for _, arch := range archiveArticles {
		var tempArch ResponseCategoryArchive
		tempArch.UUIDArticle = arch.UUIDArticle
		tempArch.URL = arch.URL
		tempArch.Header = arch.Header
		respCategoryArch = append(respCategoryArch, tempArch)
	}
	go s.Publisher(&event_bus.Event{
		Name: common.EventClickCategory,
		Data: category,
	})
	return respCategoryArch, nil
}
func (s *ServiceArticle) GetArchiveArticle(uuid string) (*model.ArticleArchive, *custom_errors.Error) {
	archArticle, errGetArchArticle := s.repo.GetArchiveArticle(uuid)
	if errGetArchArticle != nil {
		return nil, &custom_errors.Error{
			Message: ErrNotFoundArticle.Error(),
			Status:  http.StatusNotFound,
		}
	}
	go s.Publisher(&event_bus.Event{
		Name: common.EventClickArticle,
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
	tickerEveryDay := time.NewTicker(24 * time.Hour)
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
	linksList, errLinkList := s.repo.loadLinkList()
	if errLinkList != nil {
		log.Println(errLinkList)
		return
	}
	for _, list := range linksList {
		data := strings.Split(list, " ")
		var wg sync.WaitGroup
		parse := NewParsing(&wg, s.repo)
		wg.Add(3 * len(linksList))
		go parse.parseCategory(data[0], data[1])
		go parse.parseArticle(data[1])
		go parse.createRdb(data[1])
		defer func() {
			wg.Wait()
		}()
	}
}
func (s *ServiceArticle) RemoveUserArticles() {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	select {
	case <-ticker.C:

	}
}
