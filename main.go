package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/skyling/gs/alibaba"
	"github.com/tealeg/xlsx"
)

var (
	f         = flag.String("f", "list.conf", "Config File Path")
	waitgroup sync.WaitGroup
)

func init() {
	flag.Parse()
}

func main() {
	p, _ := filepath.Abs(*f)
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	var goods []*alibaba.GoodsInfo

	for {
		line, err := rd.ReadString('\n')
		if io.EOF == err {
			break
		}
		if err != nil || line == "" {
			continue
		}
		line = strings.Trim(strings.Trim(line, "\n"), "\r")
		params := strings.Split(line, " ")
		if len(params) >= 2 {
			sku, gid := params[0], params[1]
			waitgroup.Add(1)
			go func() {
				goods = append(goods, Goods(sku, gid))
				waitgroup.Done()
			}()
		}
	}
	waitgroup.Wait()
	WriteGoodsInfo(goods)
}

func WriteGoodsInfo(goods []*alibaba.GoodsInfo) {
	var (
		file       *xlsx.File
		sheet      *xlsx.Sheet
		row        *xlsx.Row
		cell       *xlsx.Cell
		saler      float64 = 1.3
		shipWeight float64 = 120
		huilv      float64 = 6
		err        error
	)
	// 分类 中文名称 英文名称 SKU 零售价(美元) (成本+挂号费+批发运费+国际运费*重量)/汇率*(销售/成本) RMB成本 挂号费 汇率 批发运费 国际运费 预估重量 销售/成本 供应商链接 特征翻译 图片翻译
	// A
	title := []string{"分类", "中文名称", "英文名称", "SKU", "零售价", "RMB成本", "挂号费", "汇率", "批发运费", "国际运费", "预估重量", "销售/成本", "供应商链接", "特征翻译"}
	var data [][]string
	data = append(data, title)
	picTrans := make(map[string]interface{})
	// fmt.Printf("%v\r\n", goods)
	if len(goods) > 0 {
		var rownum = 2
		for _, v := range goods {
			if v == nil {
				continue
			}
			guahao := float64(0)
			if (float64(v.Price)+v.TotalCost/float64(v.BeginCount)+1+shipWeight*(v.Weight+10)/1000)/huilv*saler > 5 {
				guahao = 8
			}
			featureList := ""
			if len(v.FeatureList) > 0 {
				for _, f := range v.FeatureList {
					featureList += f["name"] + " : " + f["value"] + "\r\n"
					featureList += f["enname"] + " : " + f["envalue"] + "\r\n"
				}
			}
			imgTmp := make(map[string]string)
			if len(v.DetailPics) > 0 {
				s, err := json.Marshal(v.DetailPics)
				if err == nil {
					imgTmp["detail"] = fmt.Sprintf("%s", s)
				}
			}
			if len(v.CoverPics) > 0 {
				s, err := json.Marshal(v.CoverPics)
				if err == nil {
					imgTmp["cover"] = fmt.Sprintf("%s", s)
				}
			}
			picTrans[v.SKU] = imgTmp
			tmp := []string{
				"",
				v.Names[0],
				v.Names[1],
				v.SKU,
				fmt.Sprintf("=(F%d+G%d+I%d+(J%d*K%d))/H%d*L%d", rownum, rownum, rownum, rownum, rownum, rownum, rownum),
				fmt.Sprintf("%f", v.Price),
				fmt.Sprintf("%f", guahao),
				fmt.Sprintf("%f", huilv),
				fmt.Sprintf("%f", v.TotalCost/float64(v.BeginCount)+1),
				fmt.Sprintf("%f", shipWeight),
				fmt.Sprintf("%f", (v.Weight+10)/1000),
				fmt.Sprintf("%f", saler), v.DetailURL, featureList}
			data = append(data, tmp)
			rownum++
		}
	}

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("商品信息")
	if err != nil {
		fmt.Printf(err.Error())
	}

	for _, item := range data {
		row = sheet.AddRow()
		for _, citem := range item {
			cell = row.AddCell()
			cell.Value = citem
		}
	}

	if len(picTrans) > 0 {
		for sku, pics := range picTrans {
			if pics == nil {
				continue
			}
			tmpPics := pics.(map[string]string)

			if tmpPics["detail"] == "" {
				continue
			}

			sheet, err = file.AddSheet(sku)
			if err != nil {
				fmt.Printf(err.Error())
			}

			data := jsoniter.ParseString(jsoniter.ConfigCompatibleWithStandardLibrary, tmpPics["detail"]).ReadAny()
			// map[string][]map[string]string   source_text target_text
			if len(data.Keys()) == 0 {
				continue
			}
			for _, k := range data.Keys() {
				pics := data.Get(k)
				if pics.Size() == 0 {
					continue
				}
				for i := 0; i < pics.Size(); i++ {
					row = sheet.AddRow()
					cell = row.AddCell()
					cell.Value = k
					cell = row.AddCell()
					cell.Value = pics.Get(i).Get("source_text").ToString()
					cell = row.AddCell()
					cell.Value = pics.Get(i).Get("target_text").ToString()
				}
			}
		}
	}

	err = file.Save(time.Now().Format("20060102150405") + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

// Goods 商品处理
func Goods(sku, gid string) *alibaba.GoodsInfo {
	base, _ := filepath.Abs("./resources/")
	Img := alibaba.NewImage(base, sku)
	page, err := alibaba.NewPage(sku, fmt.Sprintf("https://m.1688.com/offer/%s.html", gid))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	info := page.GetGoodsInfo()
	fmt.Println(sku + " 封面图片")
	imgURLs := page.GetCoverPics()
	img := Img.CoverPath().SetUrls(imgURLs)
	img.SaveImages()
	info.CoverPics = img.TransImages()
	detailURL := page.GetDetailURL()
	if detailURL != "" {
		detailURL = "https:" + strings.Replace(strings.Replace(detailURL, "http:", "", -1), "https:", "", -1)
		detail, err := alibaba.NewDetail(sku, detailURL)
		if err != nil {
			fmt.Println("详情获取错误!" + err.Error())
		} else {
			imgURLs := detail.GetImgs()
			fmt.Println(sku + " 详情图片")
			dimg := Img.DetailPath().SetUrls(imgURLs)
			dimg.SaveImages()
			info.DetailPics = dimg.TransImages()
		}
	}
	return info
}
