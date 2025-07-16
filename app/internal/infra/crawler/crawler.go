package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

func FetchPage(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	return goquery.NewDocumentFromReader(resp.Body)
}
