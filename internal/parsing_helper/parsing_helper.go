package parsing_helper

import (
	"app/news-parser/internal/loggers"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Browser struct {
	Logger *loggers.Logger
	*rod.Browser
	bin string
}

func NewBrowser(rodBin string, logger *loggers.Logger) *Browser {
	path, errLaunch := launcher.New().
		NoSandbox(true).
		Headless(true).
		Bin(rodBin).
		Launch()
	if errLaunch != nil {
		logger.SystemLogger(slog.LevelError, "failed to start launcher")
		os.Exit(1)
	}
	browser := rod.New().ControlURL(path).MustConnect()
	return &Browser{
		Logger:  logger,
		Browser: browser,
	}
}
func (b *Browser) RecoveryBrowser() {
	if errClose := b.Browser.Close(); errClose != nil {
		b.Logger.SystemLogger(slog.LevelError, "failed to close crashed browser")
	}
	b.Logger.SystemLogger(slog.LevelWarn, "recovery launcher")
	path, errLaunch := launcher.New().
		NoSandbox(true).
		Headless(true).
		Bin(b.bin).
		Launch()
	if errLaunch != nil {
		b.Logger.SystemLogger(slog.LevelError, "failed to start recovery launcher")
		return
	}
	b.Browser = rod.New().ControlURL(path).MustConnect()
}

func SendRequest(url string) (*http.Response, error) {
	poolUserAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x64; rv:125.0) Gecko/20100101 Firefox/125.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 OPR/108.0.0.0",
		"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.82 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (Linux; Android 14; SAMSUNG SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/24.0 Chrome/115.0.0.0 Mobile Safari/537.36",
	}
	request, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, errReq
	}
	randomIndex := rand.IntN(len(poolUserAgents))
	request.Header.Set("User-Agent", poolUserAgents[randomIndex])
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	request.Header.Set("Accept-Language", "ru,en;q=0.9,en-US;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, errResp := client.Do(request)
	if errResp != nil {
		return nil, errResp
	}
	return resp, nil
}
