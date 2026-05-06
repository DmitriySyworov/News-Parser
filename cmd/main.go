package main

import (
	"app/news-parser/configs"
	"app/news-parser/internal/article"
	"app/news-parser/internal/auth"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/open_Db"
	"app/news-parser/internal/stat"
	"app/news-parser/internal/user"
	"app/news-parser/pkg/event_bus"
	"net/http"
)

func main() {
	//conf
	conf := configs.NewConfigs()
	router := http.NewServeMux()
	//middleware
	managerMiddleware := middleware.NewManagerMiddleware(conf.Signature)
	//event_bus
	eventBus := event_bus.NewEventBus()
	//db
	redis := open_Db.OpenRedis(conf.RedisPassword, conf.RedisAddress)
	postgres := open_Db.OpenPostgres(conf.DSN)
	//repositories
	repoAuth := auth.NewRepositoryAuth(redis)
	repoUser := user.NewRepositoryUser(postgres)
	repoArticle := article.NewRepositoryArticle(postgres, redis)
	repoStat := stat.NewRepositoryStat(postgres, redis)
	//services
	serviceAuth := auth.NewServiceAuth(repoAuth, &auth.ServiceAuthDep{IRepoUser: repoUser, Configs: conf})
	serviceUser := user.NewServiceUser(repoUser)
	serviceArticle := article.NewServiceArticle(repoArticle, &article.ServiceArticleDep{EventBus: eventBus, IRepoStat: repoStat, IRepoUser: repoUser})
	serviceStat := stat.NewServiceStat(repoStat, &stat.ServiceStatDep{EventBus: eventBus})
	//goroutines
	go serviceStat.PushInStat()
	go serviceArticle.ReplacementInfo()
	//handlers
	auth.NewHandlerAuth(router, &auth.HandlerAuthDep{ServiceAuth: serviceAuth, ManagerMiddleware: managerMiddleware})
	user.NewHandlerUser(router, &user.HandlerUserDep{ServiceUser: serviceUser})
	article.NewHandlerArticle(router, &article.HandlerArticleDep{ServiceArticle: serviceArticle, ManagerMiddleware: managerMiddleware})
	stat.NewHandlerStat(router, &stat.HandlerStatDep{ServiceStat: serviceStat})
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if errApi := server.ListenAndServe(); errApi != nil {
		panic(errApi)
	}
}
