package alibaba

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/skyling/gs/requester"
)

var (
	// API_IMG_TRANS_URL 图片翻译url
	API_IMG_TRANS_URL = "http://ai.qq.com/cgi-bin/appdemo_imagetranslate"
	// API_TXT_TRANS_URL 文本翻译链接
	API_TXT_TRANS_URL = "http://ai.qq.com/cgi-bin/appdemo_texttrans"
)

// ImageTrans 图片翻译
func ImageTrans(URL string) (rets []map[string]string) {
	body, err := requester.Fetch("POST", API_IMG_TRANS_URL, map[string]string{
		"image_url": URL,
	}, map[string]string{"content-type": "application/x-www-form-urlencoded"})
	if err != nil {
		return nil
	}
	json := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, body).ReadAny()
	if json.Get("msg").ToString() != "ok" {
		return nil
	}
	records := json.Get("data").Get("image_records")
	for i := 0; i < records.Size(); i++ {
		record := records.Get(i)
		rets = append(rets, map[string]string{
			"source_text": record.Get("source_text").ToString(),
			"target_text": record.Get("target_text").ToString(),
		})
	}
	return
}

// TextTrans 文本翻译
func TextTrans(txt string) string {
	body, err := requester.Fetch("POST", API_TXT_TRANS_URL, map[string]string{
		"text": txt,
		"type": "1",
	}, map[string]string{"content-type": "application/x-www-form-urlencoded"})

	if err != nil {
		fmt.Println(err)
		return ""
	}
	json := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, body).ReadAny()
	return json.Get("data").Get("trans_text").ToString()
}
