package article_default

import (
	"app/news-parser/internal/common"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	WG        *sync.WaitGroup
	Repo      *RepositoryArticle
}

func NewParsing(wg *sync.WaitGroup, repo *RepositoryArticle) *Parse {
	return &Parse{
		LinkCh:    make(chan ArticlesGoroutines, 10),
		ArticleCh: make(chan ArticlesGoroutines, 10),
		WG:        wg,
		Repo:      repo,
	}
}

var ConditionArticles = []struct {
	Domain        string
	StartWord     string
	StopWord      string
	OpenEndpoints []string
}{
	{Domain: "https://www.bbc.com", StartWord: "Время чтения:", StopWord: "Темы", OpenEndpoints: []string{"article"}},
	//	{Domain: "https://meduza.io", StartWord: "ВОЙТИ", StopWord: "Телеграм", OpenEndpoints: []string{"feature", "news"}},
}

func (p *Parse) createRdb(category string) {
	for art := range p.ArticleCh {
		if art.Error != nil {
			log.Println(art.Error)
			continue
		}
		if art.Category == category {
			p.Repo.createNewArticle(&art)
		}
	}
	p.WG.Done()
}
func (p *Parse) parseArticle(category string) {
	path, _ := launcher.New().Headless(true).Launch()
	browser := rod.New().ControlURL(path).MustConnect()
	for art := range p.LinkCh {
		if art.Error != nil {
			p.ArticleCh <- art
		}
		defer p.recoveryGoroutine()
		ticker := time.NewTicker(30 * time.Second)
		for _, cond := range ConditionArticles {
			select {
			case <-ticker.C:
				log.Println("время на запись текста истекло")
				art.Text = "-"
				p.ArticleCh <- art
			default:
				if strings.Contains(art.Url, cond.Domain) && art.Category == category && art.IsArticle {
					defer browser.Close()
					page := browser.MustPage(art.Url)
					page.Timeout(10 * time.Second).MustWaitLoad()
					text, err := page.MustElement("body").Text()
					if err != nil {
						log.Println("проблема с парсингом страницы")
						art.Text = "-"
						p.ArticleCh <- art
					}
					res := strings.Split(text, "\n")
					startI := len(res)
					var a string
					for i := 0; i < len(res); i++ {
						if strings.Contains(res[i], cond.StartWord) {
							startI = i
						}
						if startI <= i {
							a += res[i] + "\n"
						}
						if res[i+1] == cond.StopWord {
							break
						}
					}
					art.Text = a
					p.ArticleCh <- art
				} else if !art.IsArticle {
					art.Text = "-"
					p.ArticleCh <- art
				}
			}
		}
	}
	close(p.ArticleCh)
	p.WG.Done()
}
func (p *Parse) parseCategory(url, category string) {
	for _, cond := range ConditionArticles {
		if strings.Contains(url, cond.Domain) {
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
				if linkHeader != "" && exists {
					for _, ep := range cond.OpenEndpoints {
						if strings.Contains(href, ep) {
							article.IsArticle = true
							break
						}
					}
					if !strings.Contains(href, cond.Domain) {
						article.Header = common.ParseString(linkHeader)
						article.Url = cond.Domain + common.ParseString(href)
						p.LinkCh <- article
					} else {
						article.Header = common.ParseString(linkHeader)
						article.Url = common.ParseString(href)
						p.LinkCh <- article
					}
				}
			})
		}
	}
	close(p.LinkCh)
	p.WG.Done()
}
func (p *Parse) recoveryGoroutine() {
	if errPanic := recover(); errPanic != nil {
		p.ArticleCh <- ArticlesGoroutines{Error: errors.New(fmt.Sprint(errPanic))}
	}
}
