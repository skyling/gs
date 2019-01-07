package alibaba

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/text/transform"
)

// Page 页面
type Page struct {
	SKU      string            // sku
	URL      string            // URL 地址
	ViewData jsoniter.Any      // json 格式数据
	Doc      *goquery.Document // Doc 获取到文档
	DocHtml  string            // html文档
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
func (p *Page) fetchDoc() (doc *goquery.Document, err error) {
	fmt.Println("url", p.URL)
	url := p.URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("cache-control", "no-cache")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	utfBody := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())
	doc, err = goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return
	}
	p.DocHtml, _ = doc.Html()
	p.Doc = doc
	return doc, nil
}

// GetViewData 获取商品数据
func (p *Page) GetViewData() (*Page, error) {
	// s := regexp.MustCompile(`window\.wingxViewData\[0\]=({.*}?)</script>`).FindStringSubmatch(p.DocHtml)
	// if len(s) < 2 {
	// 	return p, fmt.Errorf(p.SKU + " 详情获取失败")
	// }
	// p.ViewData = jsoniter.ParseString(jsoniter.ConfigFastest, s[1]).ReadAny()
	return p, nil
}

// GetDetailURL 获取详细信息url
func (p *Page) GetDetailURL() string {
	s := regexp.MustCompile(`"detailUrl":\s*\"(.*)?"`).FindStringSubmatch(p.DocHtml)
	// fmt.Println("detailUrl", s)
	return s[1]
}

// GetGoodsInfo 获取商品信息
func (p *Page) GetGoodsInfo() *GoodsInfo {
	goodsInfo := &GoodsInfo{}
	goodsInfo.SKU = p.SKU
	goodsInfo.DetailURL = strings.Replace(p.URL, "//m.", "//detail.", -1)
	name := p.Doc.Find("title").Text()
	goodsInfo.Names = []string{name, TextTrans(name)}
	// priceRange := p.ViewData.Get("priceRanges")
	priceRange := p.Doc.Find("#widget-wap-detail-common-price").Find("script").Text()
	// fmt.Print(priceRange)
	priceJson := jsoniter.ParseString(jsoniter.ConfigDefault, priceRange).ReadAny()
	if priceJson.Get("showPriceRanges").Size() > 0 {
		goodsInfo.Price = priceJson.Get("showPriceRanges").Get(0).Get("price").ToFloat64()
		goodsInfo.BeginCount = priceJson.Get("beginAmount").ToUint()
	}
	freightCost := p.Doc.Find("#widget-wap-detail-common-logistics > div > div.takla-item-content > span:nth-child(2) > span:nth-child(3)").Text()
	if freightCost != "" {
		goodsInfo.TotalCost, _ = strconv.ParseFloat(freightCost, 64)
	}
	freightInfo := regexp.MustCompile(`data-offer-attribute-name=\"重量\"\s*data-offer-attribute-value=\"(.*)\"`).FindStringSubmatch(p.DocHtml)
	if len(freightInfo) == 2 {
		tmp, _ := strconv.ParseFloat(freightInfo[1], 64)
		goodsInfo.Weight = tmp * 1000
	}

	p.Doc.Find(".detail-attribute-item").Each(func(i int, d *goquery.Selection) {
		name, _ := d.Attr("data-offer-attribute-name")
		value, _ := d.Attr("data-offer-attribute-value")
		goodsInfo.FeatureList = append(goodsInfo.FeatureList, map[string]string{
			"name":    name,
			"enname":  TextTrans(name),
			"value":   value,
			"envalue": TextTrans(value),
		})
	})
	return goodsInfo
}

// GetCoverPics 获取封面图片url
func (p *Page) GetCoverPics() (pics []string) {
	p.Doc.Find("#J_Detail_ImageSlidesLayer").Find("img").Each(func(i int, d *goquery.Selection) {
		src, exists := d.Attr("swipe-lazy-src")
		if exists {
			src = strings.Replace(src, "640x640", "800x800", 1)
			pics = append(pics, src)
		}
	})
	return pics
}
