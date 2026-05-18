package article_user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/model"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/google/uuid"
)

type CustomParse struct {
	ArticleUserCh chan model.UserArticle
	DBUserCh      chan model.UserArticle
	RespUserCh    chan ResponseCreateUserArticle
	IsOkTextCh    chan bool
	Repo          *RepositoryArticleUser
}

func NewCustomParsing(repo *RepositoryArticleUser) *CustomParse {
	return &CustomParse{
		ArticleUserCh: make(chan model.UserArticle, 10),
		DBUserCh:      make(chan model.UserArticle, 10),
		RespUserCh:    make(chan ResponseCreateUserArticle, 10),
		IsOkTextCh:    make(chan bool, 10),
		Repo:          repo,
	}
}
func (cp *CustomParse) customParseCategory(url, category, UserUUID string, isText bool) {
	defer cp.recoveryCustomGoroutine()
	response, errResp := http.Get(url)
	if errResp != nil {
		cp.RespUserCh <- ResponseCreateUserArticle{SuccessOperation: SuccessOfTheOperation{
			Success: false,
			Message: ErrUserURL.Error(),
			Status:  http.StatusNotFound,
		}}
		return
	}
	defer func() {
		if errClose := response.Body.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()
	doc, errParse := goquery.NewDocumentFromReader(response.Body)
	if errParse != nil {
		cp.RespUserCh <- ResponseCreateUserArticle{SuccessOperation: SuccessOfTheOperation{
			Success: false,
			Message: ErrFailedParseBody.Error(),
			Status:  http.StatusNotFound,
		}}
		return
	}
	domain := getDomain(url)
	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		var userArticle model.UserArticle
		linkHeader := element.Text()
		href, exists := element.Attr("href")
		if linkHeader != "" && exists {
			if !strings.Contains(href, domain) {
				userArticle.URL = domain + common.ParseString(href)
			} else {
				userArticle.URL = common.ParseString(href)
			}
			userArticle.ArticleUUID = uuid.New().String()
			userArticle.UserUUID = UserUUID
			userArticle.Category = category
			userArticle.Header = common.ParseString(linkHeader)

			if isText {
				cp.ArticleUserCh <- userArticle
			} else {
				cp.DBUserCh <- userArticle
			}
		}
	})
	if !isText {
		close(cp.DBUserCh)
	}
}
func getDomain(url string) string {
	counter := 0
	domain := ""
	for i := 0; i < len(url); i++ {
		if url[i] == '/' {
			counter++
		}
		if counter < 3 {
			domain += string(url[i])
		} else {
			break
		}
	}
	return domain
}

func (cp *CustomParse) CustomParseArticle() {
	path, errLaunch := launcher.New().Headless(true).Launch()
	if errLaunch != nil {
		//cp.RespUserCh <- model.UserArticle{Error: errLaunch.Error()}
		return
	}
	browser := rod.New().ControlURL(path).MustConnect()
	defer func() {
		cp.recoveryCustomGoroutine()
		if errClose := browser.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()
	for art := range cp.ArticleUserCh {
		go func() {
			page := browser.MustPage(art.URL)
			page.Timeout(10 * time.Second).MustWaitLoad()
			text, errElement := page.MustElement("body").Text()
			if errElement != nil {
				log.Println("error parse page")
				art.Text = "-"
				cp.IsOkTextCh <- true
			}
			art.Text = text
			cp.IsOkTextCh <- true
		}()
		ticker := time.NewTicker(20 * time.Second)
		select {
		case <-ticker.C:
			log.Println("time for write down text expired") //!!
			art.Text = "-"
			cp.DBUserCh <- art
		case ok := <-cp.IsOkTextCh:
			if ok {
				cp.DBUserCh <- art
			}
		}
	}
	close(cp.DBUserCh)
}
func (cp *CustomParse) createUserArticles() {
	for userArticle := range cp.DBUserCh {
		errCreate := cp.Repo.CreateUserNewArticle(&userArticle)
		if errCreate != nil {
			cp.RespUserCh <- ResponseCreateUserArticle{SuccessOperation: SuccessOfTheOperation{
				Success: false,
				Message: ErrSaveParseData.Error(),
				Status:  http.StatusInternalServerError,
			}}
			return
		}
		cp.RespUserCh <- ResponseCreateUserArticle{
			Header:      userArticle.Header,
			URL:         userArticle.URL,
			Text:        userArticle.Text,
			Category:    userArticle.Category,
			ArticleUUID: userArticle.ArticleUUID,
			UserUUID:    userArticle.UserUUID,
			SuccessOperation: SuccessOfTheOperation{
				Success: true,
				Message: "Created",
				Status:  http.StatusCreated,
			}}
	}
	close(cp.RespUserCh)
}

func (cp *CustomParse) recoveryCustomGoroutine() {
	if errPanic := recover(); errPanic != nil {
		log.Println(errPanic)
		cp.RespUserCh <- ResponseCreateUserArticle{SuccessOperation: SuccessOfTheOperation{
			Success: false,
			Message: custom_errors.ErrCriticalServer.Error(),
			Status:  http.StatusInternalServerError,
		}}
	}
}
func ParseText(url string) (string, error) {
	path, errLaunch := launcher.New().Headless(true).Launch()
	if errLaunch != nil {
		return "", errLaunch
	}
	browser := rod.New().ControlURL(path).MustConnect()
	defer func() {
		if errClose := browser.Close(); errClose != nil {
			log.Println(errClose)
		}
	}()
	page := browser.MustPage(url)
	page.Timeout(10 * time.Second).MustWaitLoad()
	text, errElement := page.MustElement("body").Text()
	if errElement != nil {
		return "", errElement
	}
	return text, nil
}
