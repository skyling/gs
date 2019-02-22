package alibaba

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/skyling/gs/requester"
)

var (
	APP_ID  = "1106705557"
	APP_KEY = "WjblYfaDw7MqLjXx"
	// API_IMG_TRANS_URL 图片翻译url
	API_IMG_TRANS_URL = "https://api.ai.qq.com/fcgi-bin/nlp/nlp_imagetranslate"
	// API_TXT_TRANS_URL 文本翻译链接
	API_TXT_TRANS_URL = "https://api.ai.qq.com/fcgi-bin/nlp/nlp_texttrans"
)

// Sign 签名
func Sign(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var str string
	for _, v := range keys {
		value := params[v]
		if value != "" {
			str = str + fmt.Sprintf("%s=%s&", v, url.QueryEscape(value))
		}
	}
	str = str + "app_key=" + APP_KEY
	sign := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(str))))
	return sign
}

// ImageTrans 图片翻译
func ImageTrans(path string) (rets []map[string]string) {
	return nil
	fbody, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	params := getParams()
	params["image"] = base64.StdEncoding.EncodeToString(fbody)
	params["session_id"] = fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(10000))
	params["scene"] = "doc"
	params["source"] = "zh"
	params["target"] = "en"
	params["sign"] = Sign(params)
	body, err := requester.Fetch("POST", API_IMG_TRANS_URL, params, map[string]string{"content-type": "application/x-www-form-urlencoded"})
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
	return ""
	params := getParams()
	params["type"] = "1"
	params["text"] = txt
	params["sign"] = Sign(params)
	body, err := requester.Fetch("POST", API_TXT_TRANS_URL, params, map[string]string{"content-type": "application/x-www-form-urlencoded"})

	if err != nil {
		fmt.Println(err)
		return ""
	}
	json := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, body).ReadAny()
	return json.Get("data").Get("trans_text").ToString()
}

func getParams() map[string]string {
	params := make(map[string]string)
	params["app_id"] = APP_ID
	params["time_stamp"] = fmt.Sprintf("%d", time.Now().Unix())
	params["nonce_str"] = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano()))))
	return params
}
