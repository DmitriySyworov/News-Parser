package main

import (
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/model"
	"app/news-parser/internal/open_Db"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	logger := loggers.NewLogger()
	godotenv.Load(".env")
	dsn := os.Getenv("DSN")
	if dsn == "" {
		logger.SystemLogger(slog.LevelError, "environment variable DSN is empty")
		os.Exit(1)
	}
	logger.SystemLogger(slog.LevelInfo, "DSN environment variable retrieved successfully")
	db := open_Db.OpenPostgres(dsn, logger)
	errMigrate := db.AutoMigrate(&model.User{}, &model.UserArticle{}, &model.UserArticleStat{}, &model.ArticleArchive{}, &model.CategoryStat{}, &model.ArticleStat{})
	if errMigrate != nil {
		logger.SystemLogger(slog.LevelError, "table migration failed")
		os.Exit(1)
	}
	logger.SystemLogger(slog.LevelInfo, "table migration successfully")
}
