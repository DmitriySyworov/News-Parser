package article_default

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/parsing_helper"
	"errors"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	Browse    *parsing_helper.Browser
	WG        *sync.WaitGroup
	Repo      *RepositoryArticle
	Logger    loggers.Logger
}

func NewParsing(wg *sync.WaitGroup, repo *RepositoryArticle, browser *parsing_helper.Browser, logger *loggers.Logger) *Parse {
	return &Parse{
		LinkCh:    make(chan ArticlesGoroutines, 10),
		ArticleCh: make(chan ArticlesGoroutines, 10),
		IsOk:      make(chan bool, 10),
		Browse:    browser,
		WG:        wg,
		Repo:      repo,
	}
}

func (p *Parse) parseCategory(url, category, domain, flagText, isArticleOnHeader string) {
	response, errResp := parsing_helper.SendRequest(url)
	if errResp != nil {
		p.LinkCh <- ArticlesGoroutines{Error: errResp}
		return
	}
	defer func() {
		if errClose := response.Body.Close(); errClose != nil {
			p.Logger.SystemLogger(slog.LevelWarn, "failed to close open link")
		}
		p.recoveryGoroutine(false)
	}()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		p.LinkCh <- ArticlesGoroutines{Error: err}
		return
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
}
func (p *Parse) parseArticle(domain, startWord, stopWord string) {
	defer p.recoveryGoroutine(false)
	for art := range p.LinkCh {
		if art.Error != nil {
			p.ArticleCh <- art
		}
		ticker := time.NewTicker(20 * time.Second)
		go func() {
			defer p.recoveryGoroutine(true)
			if strings.Contains(art.Url, domain) && art.IsArticle {
				wrapperBrowser := p.Browse.ResBrowser
				page := wrapperBrowser.MustPage(art.Url)
				page.MustWaitLoad()
				text, err := page.MustElement("body").Text()
				if err != nil {
					p.Logger.SystemLogger(slog.LevelInfo, "page parsing error")
					art.Text = "-"
					art.IsArticle = false
					if errClose := page.Close(); errClose != nil {
						p.Logger.SystemLogger(slog.LevelWarn, "failed to close page browser")
					}
					p.IsOk <- true
				}
				if errClose := page.Close(); errClose != nil {
					p.Logger.SystemLogger(slog.LevelWarn, "failed to close page browser")
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
		}()
		select {
		case <-ticker.C:
			p.Logger.SystemLogger(slog.LevelInfo, "time to write down text has expired")
			art.Text = "-"
			p.ArticleCh <- art
		case _, isOpen := <-p.IsOk:
			if isOpen {
				p.ArticleCh <- art
			}
		}
	}
}
func (p *Parse) createRdb(category string) {
	for art := range p.ArticleCh {
		if art.Error != nil {
			p.Logger.SystemLogger(slog.LevelWarn, "failed to save article in DB")
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
}
func (p *Parse) recoveryGoroutine(isPanicBrowser bool) {
	if errPanic := recover(); errPanic != nil {
		p.Logger.SystemLogger(slog.LevelWarn, "critical error while parsing links")
		if isPanicBrowser {
			p.Browse.RecoveryBrowser()
		}
		p.ArticleCh <- ArticlesGoroutines{Error: errors.New("critical error while parsing links")}
	}
}
