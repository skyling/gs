package alibaba

import (
	"fmt"
	"image"
	_ "image/gif" // 图片处理
	"image/jpeg"
	_ "image/png" // 图片处理
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

var (
	// MinWidth 最小宽度
	MinWidth = 350
	// MinHeight 最小高度
	MinHeight = 200
)

// Image 图片
type Image struct {
	Urls     []string // URLs 图片地址
	Paths    []string // Paths 保存路径
	Root     string   // Root 保存图片根路径
	Base     string   // 基础地址
	IsDetail bool     // 是否为详情图片
}

// NewImage 创建结构体
func NewImage(base, sku string) *Image {
	img := &Image{}
	img.Base = base + string(os.PathSeparator) + sku + string(os.PathSeparator)
	img.Root = img.Base
	return img
}

// CoverPath 封面图片路径
func (i *Image) CoverPath() *Image {
	i.Root = i.Base + "m"
	i.IsDetail = false
	return i
}

// DetailPath 详细图片路径
func (i *Image) DetailPath() *Image {
	i.Root = i.Base + "d"
	i.IsDetail = true
	return i
}

// SetUrls 设置图片地址
func (i *Image) SetUrls(urls []string) *Image {
	i.Urls = urls
	return i
}

// SaveImages 保存获取到的图片
func (i *Image) SaveImages() {
	saveUrls := []string{}
	for key, url := range i.Urls {
		if i.SaveImage(fmt.Sprintf("%04d", key), url) {
			saveUrls = append(saveUrls, url)
		}
	}
	i.Urls = saveUrls
}

// SaveImage 保存图片到本地
func (i *Image) SaveImage(key, url string) bool {
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("%s 获取错误\r\n", url)
		fmt.Println(err.Error())
		return false
	}
	// 创建文件
	os.MkdirAll(i.Root, os.ModePerm)
	path := strings.TrimRight(i.Root, string(os.PathSeparator)) + string(os.PathSeparator) + "pic_" + key + ".jpg"
	dst, err := os.Create(path)
	if err != nil {
		fmt.Printf("%s 创建错误 %s\r\n", url, path)
		fmt.Println(err.Error())
		return false
	}
	// 生成文件
	_, err = io.Copy(dst, res.Body)
	if err != nil {
		fmt.Printf("%s 保存错误 %s\r\n", url, path)
		fmt.Println(err.Error())
		return false
	}
	dst.Close()
	p, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	c, _, err := image.DecodeConfig(p)
	if err != nil {
		fmt.Println(err)
	}
	if c.Width < MinWidth || c.Height < MinHeight {
		fmt.Printf("%s 保存尺寸过小不保存 %d*%d\r\n", url, c.Width, c.Height)
		p.Close()
		os.Remove(path)
	} else {
		i.Paths = append(i.Paths, path)
		if i.IsDetail {
			i.ScaleImage(path, 960) // 详情图960
		} else {
			i.ScaleImage(path, 800) // 主图小于800 的
		}
		fmt.Printf("%s 保存成功 %s\r\n", url, path)
		return true
	}
	return false

}

// TransImages 翻译图片文本
func (i *Image) TransImages() map[string][]map[string]string {
	rets := make(map[string][]map[string]string)
	for key, path := range i.Paths {
		name := "pic_" + fmt.Sprintf("%04d", key)
		imgs := ImageTrans(path)
		if imgs != nil {
			rets[name] = imgs
		}
	}
	return rets
}

// ScaleImage 缩放图片
func (i *Image) ScaleImage(path string, width uint) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(width, 0, img, resize.Lanczos3)
	out, err := os.Create(filepath.Dir(path) + string(os.PathSeparator) + fmt.Sprintf("%d", width) + "_" + filepath.Base(path))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
	fmt.Println(path + " 图片大小改变成功")
	os.Remove(path)
}
