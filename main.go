package main

import (
	"bufio"
	"flag"
	"fmt"
	"gs/alibaba"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	f = flag.String("f", "list.conf", "Config File Path")
)

func init() {

}

func main() {
	flag.Parse()
	p, _ := filepath.Abs(*f)
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Trim(strings.Trim(line, "\n"), "\r")
		params := strings.Split(line, " ")
		if len(params) >= 2 {
			sku, gid := params[0], params[1]
			Goods(sku, gid)
		}
	}
}

// Goods 商品处理
func Goods(sku, gid string) {
	base, _ := filepath.Abs("./resources/")
	Img := alibaba.NewImage(base, sku)
	page, err := alibaba.NewPage(sku, fmt.Sprintf("https://m.1688.com/offer/%s.html", gid))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(sku + "封面图片")
	imgURLs := page.GetCoverPics()
	Img.CoverPath().SetUrls(imgURLs).SaveImages()

	detailURL := page.GetDetailURL()
	// fmt.Println(detailURL)
	if detailURL != "" {
		detailURL = "https:" + strings.Replace(strings.Replace(detailURL, "http:", "", -1), "https:", "", -1)
		detail, err := alibaba.NewDetail(sku, detailURL)
		if err != nil {
			fmt.Println("详情获取错误!" + err.Error())
		} else {
			imgURLs := detail.GetImgs()
			fmt.Println(sku + "详情图片")
			Img.DetailPath().SetUrls(imgURLs).SaveImages()
		}

	}
}
