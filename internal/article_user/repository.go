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
func (r *RepositoryArticleUser) UpdateUserArticle(userUUID string, userArticle *model.UserArticle) (*model.UserArticle, error) {

}
func (r *RepositoryArticleUser) GetUserArticle(userUUID string, idArticle uint) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.Where("uuid_user = ? AND id_article = ?", userUUID, idArticle).First(&userArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userArticle, nil
}
func (r *RepositoryArticleUser) DeleteAllUserArticle(uuidUser string) error {
	res := r.PostgresDb.Where("uuid_user = ?", uuidUser).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticleUser) GetAllUserArticlesWithoutText(userUUID, category string, offset, limit int) (*ResponseUserArticles, error) {
	var sliceUserArticle []model.UserArticle
	if category == "" {
		res := r.PostgresDb.
			Raw(`SELECT id, created_at, updated_at, deleted_at, header, url, category, id_article, uuid_user FROM user_articles
                    WHERE uuid_user = ?
					OFFSET ?
					LIMIT ?
`, userUUID, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res := r.PostgresDb.
			Raw(`SELECT id, created_at, updated_at, deleted_at, header, url, category, id_article, uuid_user FROM user_articles
                    WHERE uuid_user = ? AND category = ?
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
	return &ResponseUserArticles{
		SliceUserArticles: sliceUserArticle,
	}, nil
}
func (r *RepositoryArticleUser) GetAllUserArticlesWithText(userUUID, category string, offset, limit int) (*ResponseUserArticles, error) {
	var sliceUserArticle []model.UserArticle
	if category == "" {
		res := r.PostgresDb.
			Raw(`SELECT id, created_at, updated_at, deleted_at, header, url, text, category, id_article, uuid_user FROM user_articles
                    WHERE uuid_user = ?
					OFFSET ?
					LIMIT ?
`, userUUID, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res := r.PostgresDb.
			Raw(`SELECT id, created_at, updated_at, deleted_at, header, url, text, category, id_article, uuid_user FROM user_articles
                    WHERE  uuid_user = ? AND category = ?
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
	return &ResponseUserArticles{
		SliceUserArticles: sliceUserArticle,
	}, nil
}
func (r *RepositoryArticleUser) RemoveUserArticleByID(userUUID string, idArticle uint) error {
	res := r.PostgresDb.Where("uuid_user = ? AND id_article = ?", userUUID, idArticle).Delete(&model.UserArticle{})
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
	res := r.PostgresDb.Raw(`SELECT id_article, deleted_at FROM user_articles
							  WHERE deleted_at IS NOT NULL`).
		Scan(&sliceIdArticles)
	if res.Error != nil {
		log.Println(res.Error)
	}
	now := time.Now().Unix()
	for _, articles := range sliceIdArticles {
		if now-articles.DeletedAt.Unix()-common.UnixMonth > 0 {
			r.PostgresDb.Unscoped().Where("id_article = ?", articles.IdArticle).Delete(&model.UserArticle{})
		}
	}
}
func (r *RepositoryArticleUser) RecoveryUserArticle(userUUID string, idArticle int) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.
		Model(&model.UserArticle{}).
		Unscoped().
		Clauses(clause.Returning{}).
		Where("uuid_user = ? AND id_article = ? AND deleted_at IS NOT NULL", userUUID, idArticle).
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
		Where("uuid_user = ? AND deleted_at IS NOT NULL", userUUID).
		Update("deleted_at", nil).
		Scan(&sliceUserArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return sliceUserArticle, nil
}
