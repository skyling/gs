package alibaba

import (
	"fmt"
	"regexp"
	"strings"

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

// GoodsInfo 商品信息
type GoodsInfo struct {
	Names       []string            // 名称
	SKU         string              // 商品SKU
	Weight      float64             // 重量
	Price       float64             // 价格
	TotalCost   float64             // 快递费
	BeginCount  uint                // 起批个数
	FeatureList []map[string]string // 特征词
	DetailURL   string              // 地址信息
	CoverPics   map[string][]map[string]string
	DetailPics  map[string][]map[string]string
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

// GetGoodsInfo 获取商品信息
func (p *Page) GetGoodsInfo() *GoodsInfo {
	goodsInfo := &GoodsInfo{}
	goodsInfo.SKU = p.SKU
	goodsInfo.DetailURL = strings.Replace(p.URL, "//m.", "//detail.", -1)
	name := p.ViewData.Get("subject").ToString()
	goodsInfo.Names = []string{name, TextTrans(name)}
	priceRange := p.ViewData.Get("priceRanges")
	if priceRange.Size() > 0 {
		goodsInfo.Price = priceRange.Get(0).Get("price").ToFloat64()
		goodsInfo.BeginCount = priceRange.Get(0).Get("begin").ToUint()
	}
	freightCost := p.ViewData.Get("freightCost")
	if freightCost.Size() > 0 {
		goodsInfo.TotalCost = freightCost.Get(0).Get("totalCost").ToFloat64()
	}

	freightInfo := p.ViewData.Get("freightInfo")
	if freightInfo != nil {
		goodsInfo.Weight = freightInfo.Get("unitWeight").ToFloat64() * 100
	}

	featureList := p.ViewData.Get("productFeatureList")
	for i := 0; i < featureList.Size(); i++ {
		item := featureList.Get(i)
		name, value := item.Get("name").ToString(), item.Get("value").ToString()
		goodsInfo.FeatureList = append(goodsInfo.FeatureList, map[string]string{
			"name":    name,
			"enname":  TextTrans(name),
			"value":   value,
			"envalue": TextTrans(value),
		})
	}
	return goodsInfo
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
