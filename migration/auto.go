package main

import (
	"app/news-parser/internal/model"
	"app/news-parser/pkg/open_Db"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		panic(errEnv)
	}
	db := open_Db.OpenPostgres(os.Getenv("DSN"))
	errMigrate := db.AutoMigrate(&model.ArticleArchive{})
	if errMigrate != nil {
		panic(errMigrate)
	}
}
