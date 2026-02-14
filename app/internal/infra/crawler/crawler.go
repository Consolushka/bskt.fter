package crawler

import (
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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
