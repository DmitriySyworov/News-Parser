package main

import (
	"app/news-parser/configs"
	"app/news-parser/internal/article_default"
	"app/news-parser/internal/article_user"
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
	//middleware
	managerMiddleware := middleware.NewManagerMiddleware(conf.Signature)
	//router
	router := http.NewServeMux()
	//event_bus
	eventBus := event_bus.NewEventBus()
	//db
	redis := open_Db.OpenRedis(conf.RedisPassword, conf.RedisAddress)
	postgres := open_Db.OpenPostgres(conf.DSN)
	//repositories
	repoAuth := auth.NewRepositoryAuth(redis)
	repoUser := user.NewRepositoryUser(postgres, redis)
	repoArticle := article_default.NewRepositoryArticle(postgres, redis)
	repoStat := stat.NewRepositoryStat(postgres, redis)
	repoArticleUser := article_user.NewRepositoryArticleUser(postgres)
	//services
	serviceAuth := auth.NewServiceAuth(repoAuth, &auth.ServiceAuthDep{IRepoUser: repoUser, Configs: conf})
	serviceUser := user.NewServiceUser(repoUser, &user.ServiceUserDep{Configs: conf})
	serviceArticle := article_default.NewServiceArticle(repoArticle, &article_default.ServiceArticleDep{EventBus: eventBus, IRepoStat: repoStat})
	serviceStat := stat.NewServiceStat(repoStat, &stat.ServiceStatDep{EventBus: eventBus})
	serviceArticleUser := article_user.NewServiceArticleUser(repoArticleUser, &article_user.ServiceArticleUserDep{IRepoUser: repoUser})
	//goroutines
	go serviceArticle.ReplacementInfo()
	go serviceStat.PushInStat()
	go serviceArticleUser.DeletingRemoveUserArticle()
	go serviceUser.DeletingRemoveUser()
	//handlers
	auth.NewHandlerAuth(router, &auth.HandlerAuthDep{ServiceAuth: serviceAuth, ManagerMiddleware: managerMiddleware})
	user.NewHandlerUser(router, &user.HandlerUserDep{ServiceUser: serviceUser, ManagerMiddleware: managerMiddleware})
	article_default.NewHandlerArticle(router, &article_default.HandlerArticleDep{ServiceArticle: serviceArticle})
	stat.NewHandlerStat(router, &stat.HandlerStatDep{ServiceStat: serviceStat})
	article_user.NewHandlerArticleUser(router, &article_user.HandlerArticleUserDep{ServiceArticleUser: serviceArticleUser, ManagerMiddleware: managerMiddleware})
	server := http.Server{
		Addr:    ":8080",
		Handler: managerMiddleware.RecoveryPanic(router),
	}
	if errApi := server.ListenAndServe(); errApi != nil {
		panic(errApi)
	}
}
