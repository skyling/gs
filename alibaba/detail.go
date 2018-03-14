package alibaba

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Detail 详细页面
type Detail struct {
	URL string
	SKU string
	Doc *goquery.Document
}

// NewDetail 新建页面
func NewDetail(sku, url string) (*Detail, error) {
	d := &Detail{}
	d.URL = url
	d.SKU = sku
	_, err := d.fetchDoc()
	return d, err
}

// FetchDoc 从网页上获取到文档
func (p *Detail) fetchDoc() (*goquery.Document, error) {
	doc, err := goquery.NewDocument(p.URL)
	if err != nil {
		return nil, err
	}
	html, err := doc.Html()
	if err != nil {
		return nil, err
	}
	strArr := regexp.MustCompile(`var\soffer_details\=({.*});`).FindStringSubmatch(html)
	if len(strArr) < 2 {
		return nil, fmt.Errorf("数据错误")
	}
	f := strings.NewReader(strArr[1])
	p.Doc, err = goquery.NewDocumentFromReader(f)
	return p.Doc, err
}

// GetImgs 获取图片信息
func (p *Detail) GetImgs() (pics []string) {
	p.Doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imgURL, exists := s.Attr("src")
		imgURL = strings.Trim(imgURL, `\"`)
		if exists {
			pics = append(pics, imgURL)
		}
	})
	return
}
