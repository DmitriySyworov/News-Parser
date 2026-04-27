package article

import (
	"app/news-parser/internal/model"
	"app/news-parser/pkg/custom_errors"
	"strconv"
	"time"
)

type ServiceArticle struct {
	repo *RepositoryArticle
	*ServiceArticleDep
}
type ServiceArticleDep struct {
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

func NewServiceArticle(repoArticle *RepositoryArticle, dep *ServiceArticleDep) *ServiceArticle {
	return &ServiceArticle{
		repo:              repoArticle,
		ServiceArticleDep: dep,
	}
}
func (s *ServiceArticle) GetArticlesInCategoryToday(category, limitStr string) ([]ResponseCategoryToday, error) {
	if !validateCategories(category) {
		return nil, ErrCategory
	}
	limit, errParseLimit := strconv.Atoi(limitStr)
	if errParseLimit != nil {
		return nil, ErrIncorrectLimit
	}
	allArticle, errGetAllArticle := s.repo.GetArticlesInCategoryToday(category, limit)
	if errGetAllArticle != nil {
		return nil, errGetAllArticle
	}
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
	return article, nil
}
func (s *ServiceArticle) GetArticlesInCategoryArchive(category, limitStr, dateStr string) ([]ResponseCategoryArchive, error) {
	limit, errParseLimit := strconv.Atoi(limitStr)
	if errParseLimit != nil {
		return nil, ErrIncorrectLimit
	}
	date, errParseDate := time.Parse(time.DateOnly, dateStr)
	if errParseDate != nil {
		return nil, ErrIncorrectDate
	}
	archiveArticles, errGetArticlesArch := s.repo.GetArticlesInCategoryArchive(category, limit, date)
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
	return respCategoryArch, nil
}
func (s *ServiceArticle) GetArchiveArticle(uuid string) (*model.ArticleArchive, error) {
	archArticle, errGetArchArticle := s.repo.GetArchiveArticle(uuid)
	if errGetArchArticle != nil {
		return nil, custom_errors.ErrRecordNotFound
	}
	return archArticle, nil
}
func validateCategories(category string) bool {
	var StorageCategories = []string{food, politics, sport, cloth, business, electronics}
	for _, c := range StorageCategories {
		if category == c {
			return true
		}
	}
	return false
}
