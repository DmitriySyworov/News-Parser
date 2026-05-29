package main

import (
	"app/news-parser/configs"
	"app/news-parser/internal/article_default"
	"app/news-parser/internal/article_user"
	"app/news-parser/internal/auth"
	"app/news-parser/internal/event_bus"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/open_Db"
	"app/news-parser/internal/parsing_helper"
	"app/news-parser/internal/stat"
	"app/news-parser/internal/user"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	//
	logger := loggers.NewLogger()
	//
	conf := configs.NewConfigs(logger)
	//
	managerMv := middleware.NewManagerMiddleware(conf.Signature, logger)
	//
	router := http.NewServeMux()
	//
	eventBus := event_bus.NewEventBus()
	//
	browser := parsing_helper.NewBrowser(conf.RodBin, logger)
	//
	redis := open_Db.OpenRedis(conf.RedisPassword, conf.RedisAddress, logger)
	postgres := open_Db.OpenPostgres(conf.DSN, logger)
	//
	repoAuth := auth.NewRepositoryAuth(redis)
	repoUser := user.NewRepositoryUser(postgres, redis)
	repoArticle := article_default.NewRepositoryArticle(postgres, redis)
	repoStat := stat.NewRepositoryStat(postgres, redis)
	repoArticleUser := article_user.NewRepositoryArticleUser(postgres)
	//
	serviceAuth := auth.NewServiceAuth(repoAuth, &auth.ServiceAuthDep{IRepoUser: repoUser, Configs: conf})
	serviceUser := user.NewServiceUser(repoUser, &user.ServiceUserDep{Configs: conf})
	serviceArticle := article_default.NewServiceArticle(repoArticle, &article_default.ServiceArticleDep{IRepoStat: repoStat, EventBus: eventBus, Browser: browser, Logger: logger})
	serviceStat := stat.NewServiceStat(repoStat, &stat.ServiceStatDep{IRepoUser: repoUser, EventBus: eventBus})
	serviceArticleUser := article_user.NewServiceArticleUser(repoArticleUser, &article_user.ServiceArticleUserDep{IRepoUser: repoUser, EventBus: eventBus})
	//
	go serviceArticle.ReplacementInfo()
	go serviceStat.PushInStat()
	go serviceArticleUser.DeletingRemoveUserArticle()
	go serviceUser.DeletingRemoveUser()
	//
	auth.NewHandlerAuth(router, serviceAuth, &auth.HandlerAuthDep{ManagerMiddleware: managerMv, Logger: logger})
	user.NewHandlerUser(router, serviceUser, &user.HandlerUserDep{ManagerMiddleware: managerMv, Logger: logger})
	article_default.NewHandlerArticle(router, serviceArticle, &article_default.HandlerArticleDep{Logger: logger})
	stat.NewHandlerStat(router, serviceStat, &stat.HandlerStatDep{ManagerMiddleware: managerMv, Logger: logger})
	article_user.NewHandlerArticleUser(router, serviceArticleUser, &article_user.HandlerArticleUserDep{ManagerMiddleware: managerMv, Logger: logger})
	server := http.Server{
		Addr:    ":" + conf.ApiPort,
		Handler: managerMv.RecoveryPanic(managerMv.Logging(router)),
	}
	logger.SystemLogger(slog.LevelInfo, "the server started successfully on port:"+conf.ApiPort)
	if errApi := server.ListenAndServe(); errApi != nil {
		logger.SystemLogger(slog.LevelError, "critical API launch error")
		os.Exit(1)
	}
}
