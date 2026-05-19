package article_user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/model"
	"app/news-parser/internal/open_Db"
	"log"
	"time"

	"gorm.io/gorm/clause"
)

type RepositoryArticleUser struct {
	*open_Db.PostgresDb
}

func NewRepositoryArticleUser(postgres *open_Db.PostgresDb) *RepositoryArticleUser {
	return &RepositoryArticleUser{
		PostgresDb: postgres,
	}
}
func (r *RepositoryArticleUser) UpdateUserArticle(userUUID, articleUUID string, userArticle *model.UserArticle) error {
	res := r.PostgresDb.
		Where("user_uuid = ? AND article_uuid = ?", userUUID, articleUUID).
		Updates(&userArticle)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticleUser) UpdateOneColumnUserArticle(userUUID, articleUUID, nameColumn, data string) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.
		Model(&model.UserArticle{}).
		Clauses(clause.Returning{}).
		Where("user_uuid = ? AND article_uuid = ?", userUUID, articleUUID).
		Update(nameColumn, data).
		Scan(&userArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userArticle, nil
}

func (r *RepositoryArticleUser) UpdateOneColumnByDomainAll(userUUID, domain, nameColumn, data string) ([]model.UserArticle, error) {
	var userArticles []model.UserArticle
	resDomain := "%" + domain + "%"
	res := r.PostgresDb.
		Model(&model.UserArticle{}).
		Clauses(clause.Returning{}).
		Where("user_uuid = ? AND url LIKE ?", userUUID, resDomain).
		Update(nameColumn, data).
		Scan(&userArticles)
	if res.Error != nil {
		return nil, res.Error
	}
	return userArticles, nil
}
func (r *RepositoryArticleUser) GetUserArticlesByDomain(userUUID string, domain string, flagWithText bool) ([]model.UserArticle, error) {
	var sliceUserArticle []model.UserArticle
	resDomain := "%" + domain + "%"
	if flagWithText {
		res := r.PostgresDb.
			Raw(`SELECT *FROM user_articles
WHERE user_uuid = ? AND url LIKE ? AND length(text) > 1 `, userUUID, resDomain).
			Scan(&sliceUserArticle)
		if res.Error != nil || len(sliceUserArticle) == 0 {
			return nil, ErrNotFoundUserArticle
		}
	} else {
		res := r.PostgresDb.
			Raw(`SELECT *FROM user_articles
WHERE user_uuid = ? AND url LIKE ? AND length(text) <= 1 `, userUUID, resDomain).
			Scan(&sliceUserArticle)
		if res.Error != nil || len(sliceUserArticle) == 0 {
			return nil, ErrNotFoundUserArticle
		}
	}
	return sliceUserArticle, nil
}
func (r *RepositoryArticleUser) GetUserArticle(userUUID, articleUUID string) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.Where("user_uuid = ? AND article_uuid = ?", userUUID, articleUUID).First(&userArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userArticle, nil
}
func (r *RepositoryArticleUser) RemoveAllUserArticle(uuidUser string) error {
	res := r.PostgresDb.Where("user_uuid = ?", uuidUser).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticleUser) DeleteAllUserArticle(uuidUser string) error {
	res := r.PostgresDb.Unscoped().Where("user_uuid = ?", uuidUser).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticleUser) GetAllUserArticlesWithoutText(userUUID, category string, offset, limit int) (*ResponseSliceUserArticles, error) {
	var sliceUserArticle []model.UserArticle
	if category == "all" {
		res := r.PostgresDb.
			Raw(`SELECT created_at, updated_at, deleted_at, header, url, category, article_uuid, user_uuid FROM user_articles
                    WHERE user_uuid = ?  AND deleted_at IS NULL
					OFFSET ?
					LIMIT ?
`, userUUID, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res := r.PostgresDb.
			Raw(`SELECT created_at, updated_at, deleted_at, header, url, category, article_uuid, user_uuid FROM user_articles
                    WHERE user_uuid = ? AND category = ?  AND deleted_at IS NULL
					OFFSET ?
					LIMIT ?
`, userUUID, category, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	}
	if len(sliceUserArticle) == 0 {
		return nil, ErrNotFoundUserArticle
	}
	return &ResponseSliceUserArticles{
		SliceUserArticles: sliceUserArticle,
	}, nil
}
func (r *RepositoryArticleUser) GetAllUserArticlesWithText(userUUID, category string, offset, limit int) (*ResponseSliceUserArticles, error) {
	var sliceUserArticle []model.UserArticle
	if category == "all" {
		res := r.PostgresDb.
			Raw(`SELECT created_at, updated_at, deleted_at, header, url, text, category, article_uuid, user_uuid FROM user_articles
                    WHERE user_uuid = ? AND deleted_at IS NULL
					OFFSET ?
					LIMIT ?
`, userUUID, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res := r.PostgresDb.
			Raw(`SELECT created_at, updated_at, deleted_at, header, url, text, category, article_uuid, user_uuid FROM user_articles
                    WHERE  user_uuid = ? AND category = ?  AND deleted_at IS NULL
					OFFSET ?
					LIMIT ?
`, userUUID, category, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	}
	if len(sliceUserArticle) == 0 {
		return nil, ErrNotFoundUserArticle
	}
	return &ResponseSliceUserArticles{
		SliceUserArticles: sliceUserArticle,
	}, nil
}
func (r *RepositoryArticleUser) RemoveUserArticleByUUID(userUUID, articleUUID string) error {
	res := r.PostgresDb.Where("user_uuid = ? AND article_uuid = ?", userUUID, articleUUID).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticleUser) DeleteUserArticleByUUID(userUUID, articleUUID string) error {
	res := r.PostgresDb.Unscoped().Where("user_uuid = ? AND article_uuid = ?", userUUID, articleUUID).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *RepositoryArticleUser) GetRemoveUserArticle(userUUID string, offset, limit int) ([]RemoveUserArticle, error) {
	var sliceUserArticle []RemoveUserArticle
	res := r.PostgresDb.
		Model(&model.UserArticle{}).
		Unscoped().
		Where("deleted_at IS NOT NULL AND user_uuid = ?", userUUID).
		Offset(offset).
		Limit(limit).
		Scan(&sliceUserArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return sliceUserArticle, nil
}
func (r *RepositoryArticleUser) CreateUserNewArticle(art *model.UserArticle) error {
	res := r.PostgresDb.Create(&art)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

type idArticlesSlice struct {
	IdArticle uint
	DeletedAt time.Time
}

func (r *RepositoryArticleUser) deleteUserArticles() {
	var sliceIdArticles []idArticlesSlice
	res := r.PostgresDb.Raw(`SELECT article_uuid, deleted_at FROM user_articles
							  WHERE deleted_at IS NOT NULL`).
		Scan(&sliceIdArticles)
	if res.Error != nil {
		log.Println(res.Error)
	}
	now := time.Now().Unix()
	for _, articles := range sliceIdArticles {
		if now-articles.DeletedAt.Unix()-common.UnixMonth > 0 {
			r.PostgresDb.Unscoped().Where("article_uuid = ?", articles.IdArticle).Delete(&model.UserArticle{})
		}
	}
}
func (r *RepositoryArticleUser) RecoveryUserArticle(userUUID, articleUUID string) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.
		Raw(`UPDATE user_articles
				  SET deleted_at = null
				  WHERE user_uuid = ? AND article_uuid = ? AND deleted_at IS NOT NULL
				  RETURNING created_at, updated_at, deleted_at, header, url, text, category, article_uuid, user_uuid`,
			userUUID, articleUUID).
		Scan(&userArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userArticle, nil
}
func (r *RepositoryArticleUser) RecoveryAllUserArticle(userUUID string) ([]model.UserArticle, error) {
	var sliceUserArticle []model.UserArticle
	res := r.PostgresDb.
		Raw(`UPDATE user_articles
			     SET deleted_at = null
			     WHERE user_uuid = ? AND deleted_at IS NOT NULL
			     RETURNING created_at, updated_at, deleted_at, header, url, text, category, article_uuid, user_uuid`,
			userUUID).
		Scan(&sliceUserArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return sliceUserArticle, nil
}
func (r *RepositoryArticleUser) IsDomainArticleExist(userUUID, domain string) bool {
	resDomain := "%" + domain + "%"
	res := r.PostgresDb.
		Where("user_uuid = ? AND url LIKE ?", userUUID, resDomain).
		First(&model.UserArticle{})
	if res.Error != nil {
		return false
	}
	return true
}
func (r *RepositoryArticleUser) IsUserArticleRemoveExistByUUID(userUUID, articleUUID string) bool {
	res := r.PostgresDb.Unscoped().
		Where("user_uuid = ? AND article_uuid = ? AND deleted_at IS NOT NULL", userUUID, articleUUID).
		First(&model.UserArticle{})
	if res.Error != nil {
		return false
	}
	return true
}
func (r *RepositoryArticleUser) IsUserArticleExistNoRemoveAll(userUUID string) bool {
	res := r.PostgresDb.
		Where("user_uuid = ?", userUUID).
		First(&model.UserArticle{})
	if res.Error != nil {
		return false
	}
	return true
}
func (r *RepositoryArticleUser) IsUserArticleExistAll(userUUID string) bool {
	res := r.PostgresDb.
		Unscoped().
		Where("user_uuid = ?", userUUID).
		First(&model.UserArticle{})
	if res.Error != nil {
		return false
	}
	return true
}
func (r *RepositoryArticleUser) IsUserArticleRecoveryExist(userUUID string) bool {
	res := r.PostgresDb.Unscoped().
		Unscoped().
		Where("user_uuid = ? AND deleted_at IS NOT NULL", userUUID).
		First(&model.UserArticle{})
	if res.Error != nil {
		return false
	}
	return true
}
func (r *RepositoryArticleUser) IsUserArticleExistByUUID(userUUID, articleUUID string) bool {
	res := r.PostgresDb.Unscoped().
		Unscoped().
		Where("user_uuid = ? AND article_uuid = ?", userUUID, articleUUID).
		First(&model.UserArticle{})
	if res.Error != nil {
		return false
	}
	return true
}

func (r *RepositoryArticleUser) IsUserArticleExist(userUUID, articleUUID string) bool {
	res := r.PostgresDb.
		Where("user_uuid = ? AND article_uuid = ? ", userUUID, articleUUID).
		First(&model.UserArticle{})
	if res.Error != nil {
		return false
	}
	return true
}
