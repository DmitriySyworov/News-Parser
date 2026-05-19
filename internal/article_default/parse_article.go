package article_default

import (
	"app/news-parser/internal/common"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type ArticlesGoroutines struct {
	Header    string
	Url       string
	IsArticle bool
	Text      string
	Category  string
	Error     error
}
type Parse struct {
	LinkCh    chan ArticlesGoroutines
	ArticleCh chan ArticlesGoroutines
	IsOk      chan bool
	//Timeout   context.Context
	WG   *sync.WaitGroup
	Repo *RepositoryArticle
}

func NewParsing(wg *sync.WaitGroup, repo *RepositoryArticle) *Parse { //, timeout context.Context) *Parse {
	return &Parse{
		LinkCh:    make(chan ArticlesGoroutines, 10),
		ArticleCh: make(chan ArticlesGoroutines, 10),
		IsOk:      make(chan bool, 10),
		WG:        wg,
		//Timeout:   timeout,
		Repo: repo,
	}
}

func (p *Parse) parseCategory(url, category, domain, flagText, isArticleOnHeader string) {
	//go func() {
	response, errResp := http.Get(url)
	if errResp != nil {
		p.LinkCh <- ArticlesGoroutines{Error: errResp}
	}
	defer func() {
		if errClose := response.Body.Close(); errClose != nil {
			fmt.Println(errClose)
		}
		p.recoveryGoroutine()
	}()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		p.LinkCh <- ArticlesGoroutines{Error: errResp}
	}
	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		var article ArticlesGoroutines
		article.Category = category
		linkHeader := element.Text()
		href, exists := element.Attr("href")
		parsHeader := common.ParseString(linkHeader)
		if parsHeader != "" && exists && len(parsHeader) < 100 {
			if ok, _ := regexp.Match(flagText, []byte(href)); ok {
				article.IsArticle = true
			}
			if strings.Contains(parsHeader, isArticleOnHeader) {
				article.IsArticle = false
			}
			if !strings.Contains(href, domain) {
				article.Header = parsHeader
				article.Url = domain + common.ParseString(href)
				p.LinkCh <- article
			} else {
				article.Header = parsHeader
				article.Url = common.ParseString(href)
				p.LinkCh <- article
			}
		}
	})
	//	select {
	//	case <-p.Timeout.Done():
	//		return
	//	}
	//}()
	//select {
	//case <-p.Timeout.Done():
	//	return
	//}
}
func (p *Parse) parseArticle(domain, startWord, stopWord string) {
	//defer p.recoveryGoroutine()
	//go func() {
	path, errLaunch := launcher.New().Headless(true).Launch()
	if errLaunch != nil {
		log.Println(errLaunch)
	}
	browser := rod.New().ControlURL(path).MustConnect()
	defer func() {
		if errClose := browser.Close(); errClose != nil {
			log.Println(errClose)
		}
		p.recoveryGoroutine()
	}()
	for art := range p.LinkCh {
		if art.Error != nil {
			p.ArticleCh <- art
		}
		defer p.recoveryGoroutine()
		ticker := time.NewTicker(20 * time.Second)
		go func() {
			defer p.recoveryGoroutine()
			if strings.Contains(art.Url, domain) && art.IsArticle {
				page := browser.MustPage(art.Url)
				page.Timeout(10 * time.Second).MustWaitLoad()
				text, err := page.MustElement("body").Text()
				if err != nil {
					log.Println("page parsing error")
					art.Text = "-"
					art.IsArticle = false
					if errClose := page.Close(); errClose != nil {
						log.Println(errClose)
					}
					p.IsOk <- true
				}
				if errClose := page.Close(); errClose != nil {
					log.Println(errClose)
				}
				res := strings.Split(text, "\n")
				startI := len(res)
				var a string
				for i := 0; i < len(res); i++ {
					if strings.Contains(res[i], startWord) {
						startI = i
					}
					if startI <= i {
						a += res[i] + "\n"
					}
					if res[i+1] == stopWord {
						break
					}
				}
				art.Text = a
				p.IsOk <- true
			} else if !art.IsArticle {
				art.IsArticle = false
				art.Text = "-"
				p.IsOk <- true

			}
			//select {
			//case <-p.Timeout.Done():
			//	return
			//}
		}()
		select {
		//case <-p.Timeout.Done():
		//	return
		case <-ticker.C:
			log.Println("time to write down text has expired")
			art.Text = "-"
			p.ArticleCh <- art
		case _, isOpen := <-p.IsOk:
			if isOpen {
				p.ArticleCh <- art
			}
		}
	}
	//	select {
	//	case <-p.Timeout.Done():
	//		return
	//	}
	//}()
	//select {
	//case <-p.Timeout.Done():
	//	return
	//}
}
func (p *Parse) createRdb(category string) {
	//go func() {
	for art := range p.ArticleCh {
		if art.Error != nil {
			log.Println(art.Error)
			continue
		}
		if art.Text == "" || art.Text == "-" {
			art.IsArticle = false
		} else {
			art.IsArticle = true
		}
		if art.Category == category {
			p.Repo.createNewArticle(&art)
		}
	}
	//select {
	//case <-p.Timeout.Done():
	//	return
	//}
	//}()
	//select {
	//case <-p.Timeout.Done():
	//	return
	//}
}
func (p *Parse) recoveryGoroutine() {
	if errPanic := recover(); errPanic != nil {
		log.Println(errPanic)
		p.ArticleCh <- ArticlesGoroutines{Error: errors.New(fmt.Sprint(errPanic))}
	}
}
