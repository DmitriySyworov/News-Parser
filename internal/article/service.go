package article

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/model"
	"app/news-parser/pkg/event_bus"
	"log"
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
	di.IRepoUser
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
func (s *ServiceArticle) GetArticlesInCategoryToday(category, limitStr, filter string) ([]ResponseCategoryToday, error) {
	if !validateCategories(category) {
		return nil, ErrCategory
	}
	limit, errParseLimit := strconv.Atoi(limitStr)
	if errParseLimit != nil {
		return nil, ErrIncorrectParams
	}
	allArticle, errGetAllArticle := s.repo.GetArticlesInCategoryToday(category, filter, limit)
	if errGetAllArticle != nil {
		return nil, errGetAllArticle
	}
	go s.Publisher(&event_bus.Event{
		Name: common.EventClickCategory,
		Data: category,
	})
	return allArticle, nil
}
func (s *ServiceArticle) GetArticleToday(idStr string) (*model.ArticleToday, error) {
	id, errParseId := strconv.Atoi(idStr)
	if errParseId != nil {
		return nil, ErrIncorrectId
	}
	article, errGetArticle := s.repo.GetArticleToday(id)
	if errGetArticle != nil {
		return nil, errGetArticle
	}
	go s.Publisher(&event_bus.Event{
		Name: common.EventClickArticle,
		Data: article.URL,
	})
	return article, nil

}
func (s *ServiceArticle) GetArticlesInCategoryArchive(category, offsetStr, limitStr, dateStr string) ([]ResponseCategoryArchive, error) {
	if !validateCategories(category) {
		return nil, ErrCategory
	}
	limit, errParse := strconv.Atoi(limitStr)
	if errParse != nil {
		return nil, ErrIncorrectParams
	}
	var offset int
	if offsetStr == "" {
		offset = 0
	} else {
		offset, errParse = strconv.Atoi(offsetStr)
		if errParse != nil {
			return nil, ErrIncorrectParams
		}
	}
	date, errParseDate := time.Parse(time.DateOnly, dateStr)
	if errParseDate != nil {
		return nil, custom_errors.ErrIncorrectDate
	}
	archiveArticles, errGetArticlesArch := s.repo.GetArticlesInCategoryArchive(category, offset, limit, date)
	if errGetArticlesArch != nil {
		return nil, errGetArticlesArch
	}
	var respCategoryArch []ResponseCategoryArchive
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
func (s *ServiceArticle) GetArchiveArticle(uuid string) (*model.ArticleArchive, error) {
	archArticle, errGetArchArticle := s.repo.GetArchiveArticle(uuid)
	if errGetArchArticle != nil {
		return nil, custom_errors.ErrRecordNotFound
	}
	go s.Publisher(&event_bus.Event{
		Name: common.EventClickArticle,
		Data: archArticle.URL,
	})
	return archArticle, nil
}
func (s *ServiceArticle) CreateUserArticles(body *RequestCreateArticle, uuid string, isAddText bool) (*ResponseUserArticles, error) {
	if _, errNotFound := s.IRepoUser.GetUserByUUID(uuid); errNotFound != nil {
		return nil, custom_errors.ErrUserNotFound
	}
	var sliceUserArticle []model.UserArticle
	var wg sync.WaitGroup
	customParsing := NewCustomParsing(&wg, s.repo)
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
func (s *ServiceArticle) DeleteUserArticle(idArticleStr, uuidUser, allArticleStr string) error {
	var isAllArticle bool
	if allArticleStr == "true" {
		isAllArticle = true
	} else if allArticleStr == "false" {
		isAllArticle = false
	}
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		return ErrIncorrectId
	}
	if isAllArticle {
		if errDeleteAll := s.repo.DeleteAllUserArticle(uuidUser); errDeleteAll != nil {
			return custom_errors.ErrRecordNotFound
		}
	} else {
		if errDelete := s.repo.DeleteUserArticleByID(uint(idArticle)); errDelete != nil {
			return custom_errors.ErrRecordNotFound
		}
	}
	return nil
}
func (s *ServiceArticle) GetUserArticle(idArticleStr string) (*model.UserArticle, error) {
	idArticle, errParseId := strconv.Atoi(idArticleStr)
	if errParseId != nil {
		return nil, ErrIncorrectId
	}
	if userArticle, errGetUserArt := s.repo.GetUserArticle(uint(idArticle)); errGetUserArt != nil {
		return nil, custom_errors.ErrUserNotFound
	} else {
		return userArticle, nil
	}
}
func (s *ServiceArticle) GetAllUserArticles(category, offsetStr, limitStr, withTextStr string) (*ResponseUserArticles, error) {
	limit, errParse := strconv.Atoi(limitStr)
	if errParse != nil {
		return nil, ErrIncorrectParams
	}
	var offset int
	if offsetStr == "" {
		offset = 0
	} else {
		offset, errParse = strconv.Atoi(offsetStr)
		if errParse != nil {
			return nil, ErrIncorrectParams
		}
	}
	var withText bool
	if withTextStr == "true" {
		withText = true
	} else if withTextStr == "false" {
		withText = false
	} else {
		return nil, ErrIncorrectParams
	}
	if withText {
		respUserArticles, errGetArticles := s.repo.GetAllUserArticlesWithText(category, offset, limit)
		if errGetArticles != nil {
			return nil, custom_errors.ErrRecordNotFound
		}
		return respUserArticles, nil
	} else {
		respUserArticles, errGetArticles := s.repo.GetAllUserArticlesWithoutText(category, offset, limit)
		if errGetArticles != nil {
			return nil, custom_errors.ErrRecordNotFound
		}
		return respUserArticles, nil
	}
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
	tickerCheckEmpty := time.NewTicker(1 * time.Minute) //!!!
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
