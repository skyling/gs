package alibaba

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
)

// Page 页面
type Page struct {
	SKU      string            // sku
	URL      string            // URL 地址
	ViewData jsoniter.Any      // json 格式数据
	Doc      *goquery.Document // Doc 获取到文档
}

// NewPage 新建页面
func NewPage(sku, url string) (*Page, error) {
	p := &Page{}
	p.URL = url
	p.SKU = sku
	_, err := p.fetchDoc()
	if err != nil {
		return nil, err
	}
	return p.GetViewData()
}

// FetchDoc 从网页上获取到文档
func (p *Page) fetchDoc() (*goquery.Document, error) {
	doc, err := goquery.NewDocument(p.URL)
	if err != nil {
		return nil, err
	}
	p.Doc = doc
	return doc, nil
}

// GetViewData 获取商品数据
func (p *Page) GetViewData() (*Page, error) {
	html, _ := p.Doc.Html()
	s := regexp.MustCompile(`window\.wingxViewData\[0\]=({.*}?)</script>`).FindStringSubmatch(html)
	if len(s) < 2 {
		return p, fmt.Errorf(p.SKU + " 详情获取失败")
	}
	p.ViewData = jsoniter.ParseString(jsoniter.ConfigFastest, s[1]).ReadAny()
	return p, nil
}

// GetDetailURL 获取详细信息url
func (p *Page) GetDetailURL() string {
	return p.ViewData.Get("detailUrl").ToString()
}

// GetCoverPics 获取封面图片url
func (p *Page) GetCoverPics() (pics []string) {
	imgs := p.ViewData.Get("imageList")

	for i := 0; i < imgs.Size(); i++ {
		img := imgs.Get(i)
		uri := img.Get("originalImageURI").ToString()
		if uri != "" {
			pics = append(pics, uri)
		}
	}
	return pics
}
