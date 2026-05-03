package main

import (
	"app/news-parser/configs"
	"app/news-parser/internal/article"
	"app/news-parser/internal/open_Db"
	"app/news-parser/internal/stat"
	"app/news-parser/pkg/event_bus"
	"net/http"
)

func main() {
	//conf
	conf := configs.NewConfigs()
	router := http.NewServeMux()
	//event_bus
	eventBus := event_bus.NewEventBus()
	//db
	redis := open_Db.OpenRedis(conf.RedisPassword, conf.RedisAddress)
	postgres := open_Db.OpenPostgres(conf.DSN)
	//repositories
	repoArticle := article.NewRepositoryArticle(postgres, redis)
	repoStat := stat.NewRepositoryStat(postgres, redis)
	//services
	serviceArticle := article.NewServiceArticle(repoArticle, &article.ServiceArticleDep{EventBus: eventBus, IRepoStat: repoStat})
	serviceStat := stat.NewServiceStat(repoStat, &stat.ServiceStatDep{EventBus: eventBus})
	//goroutines
	go serviceStat.PushInStat()
	go serviceArticle.ReplacementInfo()
	//handlers
	article.NewHandlerArticle(router, &article.HandlerArticleDep{ServiceArticle: serviceArticle})
	stat.NewHandlerStat(router, &stat.HandlerStatDep{ServiceStat: serviceStat})
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if errApi := server.ListenAndServe(); errApi != nil {
		panic(errApi)
	}
}
