package sample_1

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sync"
	"sync/atomic"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/8/13 16:06
@description:
**/

func TestCompareAndSwap(t *testing.T) {
	var old int32 = 2

	f := func() {
		if atomic.CompareAndSwapInt32(&old, 2, 3) {
			fmt.Println("swap successful")
		} else {
			fmt.Println("failure to swap")
		}
	}
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}
	wg.Wait()

}

func TestSignture(t *testing.T) {
	data := struct {
		ProviderIds    []int  `json:"providerIds"`
		Env            string `json:"env"`
		TaskId         string `json:"taskId"`
		PackageName    string `json:"packageName"`
		AppName        string `json:"appName"`
		CloudType      int    `json:"cloudType"`
		PackageVersion string `json:"packageVersion"`
		URL            string `json:"url"`
		URLMd5         string `json:"urlMd5"`
		GameId         *int   `json:"gameId"` // 使用指针来表示可能的nil值
		Version        string `json:"version"`
		Action         string `json:"action"`
	}{
		ProviderIds: []int{10201},
		PackageName: "ctyTest",
		CloudType:   0,
		Env:         "STAGING",
		TaskId:      "1825374018769207296",
	}

	// 将数据序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	// 创建一个bytes.Buffer来存放JSON数据
	var jsonBuffer bytes.Buffer
	jsonBuffer.Write(jsonData)

	// 计算BodyMd5
	hash := md5.New()
	hash.Write(jsonBuffer.Bytes())
	bodyMd5 := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	// 打印JSON字符串
	fmt.Println(bodyMd5)

	hmac := hmac.New(sha1.New, []byte("48214359b96bf914a2de9d7bd609070f"+"&"))
	hmac.Write([]byte("POST&/cloud/package/addGame&OsU/PS9uC+wFO6Zn0U6cDQ==&AccessKey=8581e340435757dcd6ba638e358640b5&Format=JSON&SignatureMethod=HMAC-SHA1&SignatureNonce=824282677&SignatureVersion=1.1&Timestamp=1723779616782&Version=2020-11-01"))
	lsignature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	fmt.Println(lsignature)
}

func TestRegex(t *testing.T) {
	var originURLRegex *regexp.Regexp
	sub := regexp.MustCompile(`^r\((.+)\)$`).FindStringSubmatch("/cloud/package/progressQuery")
	if len(sub) == 2 {
		originURLRegex = regexp.MustCompile(sub[1])
	}

	u := "http://beta-06-api-bnd.yuntiancloud.com:18174/cloud/package/progressQuery?AccessKey=8581e340435757dcd6ba638e358640b5&Format=JSON&SignatureMethod=HMAC-SHA1&SignatureNonce=362665173&SignatureVersion=1.1&Timestamp=1724055129368&Version=2020-11-01&Signature=w2VuIlizp6qkwJl7gzcwPVh4/c0="
	up, _ := url.Parse(u)
	tr := originURLRegex.FindSubmatch([]byte(up.RequestURI()))
	fmt.Println(tr)

}
