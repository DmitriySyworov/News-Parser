package main

import (
	"app/news-parser/configs"
	"app/news-parser/internal/article"
	"app/news-parser/pkg/open_Db"
	"net/http"
)

func main() {
	//conf
	conf := configs.NewConfigs()
	router := http.NewServeMux()
	//db
	redis := open_Db.OpenRedis(conf.RedisPassword, conf.RedisAddress)
	postgres := open_Db.OpenPostgres(conf.DSN)
	//repositories
	repoArticle := article.NewRepositoryArticle(postgres, redis)
	//services
	serviceArticle := article.NewServiceArticle(repoArticle, &article.ServiceArticleDep{})
	//handlers
	article.NewHandlerArticle(router, &article.HandlerArticleDep{ServiceArticle: serviceArticle})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if errApi := server.ListenAndServe(); errApi != nil {
		panic(errApi)
	}
}
