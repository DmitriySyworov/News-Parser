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
func (r *RepositoryArticleUser) UpdateUserArticle(userUUID string, userArticle *model.UserArticle) error {
	res := r.PostgresDb.Where("user_uuid = ?", userUUID).Updates(&userArticle)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticleUser) UpdateOneColumnUserArticle(userUUID string, data string, nameColumn string) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.
		Model(&model.UserArticle{}).
		Clauses(clause.Returning{}).
		Where("user_uuid = ?", userUUID).
		Update(nameColumn, data).
		Scan(&userArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userArticle, nil
}
func (r *RepositoryArticleUser) GetUserArticlesByDomain(userUUID string, domain string, flagWithText bool) ([]model.UserArticle, error) {
	var sliceUserArticle []model.UserArticle
	if flagWithText {
		res := r.PostgresDb.Where("user_uuid = ? AND url LIKE %?% AND text IS NOT NULL OR text != '-'", userUUID, domain).First(&sliceUserArticle)
		if res.Error != nil || len(sliceUserArticle) == 0 {
			return nil, ErrNotFoundUserArticle
		}
	} else {
		res := r.PostgresDb.Where("user_uuid = ? AND url LIKE %?% AND text IS NULL OR text = '-'", userUUID, domain).First(&sliceUserArticle)
		if res.Error != nil || len(sliceUserArticle) == 0 {
			return nil, ErrNotFoundUserArticle
		}
	}
	return sliceUserArticle, nil
}
func (r *RepositoryArticleUser) GetUserArticle(userUUID string, idArticle uint) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.Where("user_uuid = ? AND article_uuid = ?", userUUID, idArticle).First(&userArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userArticle, nil
}
func (r *RepositoryArticleUser) DeleteAllUserArticle(uuidUser string) error {
	res := r.PostgresDb.Where("user_uuid = ?", uuidUser).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticleUser) GetAllUserArticlesWithoutText(userUUID, category string, offset, limit int) (*ResponseSliceUserArticles, error) {
	var sliceUserArticle []model.UserArticle
	if category == "" {
		res := r.PostgresDb.
			Raw(`SELECT created_at, updated_at, deleted_at, header, url, category, article_uuid, user_uuid FROM user_articles
                    WHERE user_uuid = ?
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
                    WHERE user_uuid = ? AND category = ?
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
	if category == "" {
		res := r.PostgresDb.
			Raw(`SELECT created_at, updated_at, deleted_at, header, url, text, category, article_uuid, user_uuid FROM user_articles
                    WHERE user_uuid = ?
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
                    WHERE  user_uuid = ? AND category = ?
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
func (r *RepositoryArticleUser) RemoveUserArticleByID(userUUID string, idArticle uint) error {
	res := r.PostgresDb.Where("user_uuid = ? AND article_uuid = ?", userUUID, idArticle).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *RepositoryArticleUser) GetRemoveUserArticle(uuid string, offset, limit int) ([]RemoveUserArticle, error) {
	var sliceUserArticle []RemoveUserArticle
	res := r.PostgresDb.
		Model(&model.UserArticle{}).
		Unscoped().
		Where("deleted_at IS NOT NULL AND uuid = ?", uuid).
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
func (r *RepositoryArticleUser) RecoveryUserArticle(userUUID string, idArticle int) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.
		Model(&model.UserArticle{}).
		Unscoped().
		Clauses(clause.Returning{}).
		Where("user_uuid = ? AND article_uuid = ? AND deleted_at IS NOT NULL", userUUID, idArticle).
		Update("deleted_at", nil).
		Scan(&userArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userArticle, nil
}
func (r *RepositoryArticleUser) RecoveryAllUserArticle(userUUID string) ([]model.UserArticle, error) {
	var sliceUserArticle []model.UserArticle
	res := r.PostgresDb.
		Model(&model.UserArticle{}).
		Unscoped().
		Clauses(clause.Returning{}).
		Where("user_uuid = ? AND deleted_at IS NOT NULL", userUUID).
		Update("deleted_at", nil).
		Scan(&sliceUserArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return sliceUserArticle, nil
}
