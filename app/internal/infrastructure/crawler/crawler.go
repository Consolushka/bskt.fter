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

	//var values []string
	//doc.Find(".PlayerSummary_playerInfoValue__JS8_v").Each(func(i int, s *goquery.Selection) {
	//	values = append(values, s.Text())
	//})
}
