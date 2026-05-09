package article_user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/model"
	"app/news-parser/pkg/generate_random"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type CustomParse struct {
	LinkUserCh    chan model.UserArticle
	ArticleUserCh chan model.UserArticle
	IsOkTextCh    chan bool
	WG            *sync.WaitGroup
	Repo          *RepositoryArticleUser
}

func NewCustomParsing(wg *sync.WaitGroup, repo *RepositoryArticleUser) *CustomParse {
	return &CustomParse{
		LinkUserCh:    make(chan model.UserArticle, 10),
		ArticleUserCh: make(chan model.UserArticle, 10),
		IsOkTextCh:    make(chan bool, 10),
		WG:            wg,
		Repo:          repo,
	}
}
func (cp *CustomParse) customParseCategory(url, category, uuid string, isText bool) {
	defer cp.recoveryCustomGoroutine()
	after := time.After(time.Second * 10)
	select {
	case <-after:

	}
	response, errResp := http.Get(url)
	if errResp != nil {
		cp.LinkUserCh <- model.UserArticle{Error: errResp.Error()}
	}
	defer func() {
		if errClose := response.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()
	doc, errParse := goquery.NewDocumentFromReader(response.Body)
	if errParse != nil {
		cp.LinkUserCh <- model.UserArticle{Error: errParse.Error()}
	}
	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		var userArticle model.UserArticle
		linkHeader := element.Text()
		href, exists := element.Attr("href")
		if linkHeader != "" && exists {
			userArticle.IDArticle = uint(generate_random.GenerateNumbers(11))
			userArticle.UUIDUser = uuid
			userArticle.Category = category
			userArticle.URL = common.ParseString(href)
			userArticle.Header = common.ParseString(linkHeader)
			cp.LinkUserCh <- userArticle
		}
	})
	cp.WG.Done()
	if isText {
		close(cp.LinkUserCh)
	}
}
func (cp *CustomParse) CustomParseArticle() {
	path, errLaunch := launcher.New().Headless(true).Launch()
	if errLaunch != nil {
		cp.ArticleUserCh <- model.UserArticle{Error: errLaunch.Error()}
	}
	browser := rod.New().ControlURL(path).MustConnect()
	defer func() {
		cp.recoveryCustomGoroutine()
		if errClose := browser.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()
	for art := range cp.LinkUserCh {
		if art.Error != "" {
			cp.ArticleUserCh <- art
		}
		go func() {
			page := browser.MustPage(art.URL)
			page.Timeout(10 * time.Second).MustWaitLoad()
			text, errElement := page.MustElement("body").Text()
			if errElement != nil {
				log.Println("проблема с парсингом страницы") //!!
				art.Text = "-"
				cp.IsOkTextCh <- true
			}
			art.Text = text
			cp.IsOkTextCh <- true
		}()
		ticker := time.NewTicker(20 * time.Second)
		select {
		case <-ticker.C:
			log.Println("время на запись текста истекло") //!!
			art.Text = "-"
			cp.ArticleUserCh <- art
		case ok := <-cp.IsOkTextCh:
			if ok {
				cp.ArticleUserCh <- art
			}
		}
	}
	cp.WG.Done()
}
func (cp *CustomParse) createUserArticlesWithoutText() {
	for userArticle := range cp.LinkUserCh {
		errCreate := cp.Repo.CreateUserNewArticle(&userArticle)
		if errCreate != nil {
			cp.LinkUserCh <- model.UserArticle{URL: userArticle.URL, Error: errCreate.Error()}
		}
		cp.LinkUserCh <- userArticle
	}
	cp.WG.Done()
	close(cp.LinkUserCh)
}
func (cp *CustomParse) createUserArticlesWithText() {
	for userArticle := range cp.ArticleUserCh {
		errCreate := cp.Repo.CreateUserNewArticle(&userArticle)
		if errCreate != nil {
			cp.ArticleUserCh <- model.UserArticle{URL: userArticle.URL, Error: errCreate.Error()}
		}
		cp.ArticleUserCh <- userArticle
	}
	cp.WG.Done()
	close(cp.ArticleUserCh)
}
func (cp *CustomParse) recoveryCustomGoroutine() {
	if errPanic := recover(); errPanic != nil {
		cp.ArticleUserCh <- model.UserArticle{Error: fmt.Sprint(errPanic)}
	}
}
