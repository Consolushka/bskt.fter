package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func FetchPage(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return goquery.NewDocumentFromReader(resp.Body)
}
