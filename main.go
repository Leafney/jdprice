package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.65 Safari/537.36"

// getGoodsApiUrl 拼接商品的请求url链接
// skuids 单个商品编码或多个商品编码的集合
func getGoodsApiUrl(skuids interface{}) string {
	// url = "https://p.3.cn/prices/mgets?callback=jQuery8844217&type=1&area=1_72_2799_0&pdtk=&pduid=336440463&pdpin=&pin=null&pdbp=0&skuIds=J_1580032&ext=11000000&source=item-pc"
	base_url := "https://p.3.cn/prices/mgets"

	// 通过断言确定skuids是一个商品的编码string类型还是多个商品的编码[]string 类型
	var paramSlice []string
	switch sku := skuids.(type) {
	case string:
		paramSlice = append(paramSlice, sku)
	case []string:
		paramSlice = sku
	default:
		panic("skuids type not supported")
	}
	// fmt.Println(paramSlice)

	params := make(url.Values)
	params.Set("type", "1")
	params.Set("source", "item-pc")
	params.Set("pduid", "336440463") // Todo:获取随机值

	// 将多个商品的编码组合成逗号分隔的字符串
	// skuidStr := strings.Join(paramSlice, ",")
	// fmt.Println(skuidStr)
	params.Set("skuIds", strings.Join(paramSlice, ","))

	//Todo: 省_市_区_县
	params.Set("area", "1_72_2799_0")

	api_url := base_url + "?" + params.Encode()
	return api_url
}

// getPriceOnly 获取价格
func getPriceOnly(skuids interface{}) {
	// 拼接请求Url
	reqUrl := getGoodsApiUrl(skuids)
	jsonData := getHttpReq(reqUrl)
	//输出返回结果
	// fmt.Println(string(jsonData))

	var jdPrice Items
	err := json.Unmarshal(jsonData, &jdPrice)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(jdPrice)
	// fmt.Println(jdPrice[0].P)

	for _, v := range jdPrice {
		fmt.Println("商品价格:", v.P)
	}

}

type Items []Priceinfo

type Priceinfo struct {
	P  string `json:"p"`
	OP string `json:"op"`
	Id string `json:"id"`
	M  string `json:"m"`
}

// 发送请求
func getHttpReq(reqUrl string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/x-javascript; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// respBody := string(body)
	// fmt.Println(respBody)
	// return respBody

	// 返回 []byte
	return body
}

func main() {

	// 拼接单个商品的价格链接
	// x := getGoodsApiUrl("2720472")
	// fmt.Println(x)

	// // 拼接多个商品的价格链接
	// goods := []string{"2720472", "1580032"}
	// y := getGoodsApiUrl(goods)
	// fmt.Println(y)

	//获取单个商品价格
	// getPriceOnly("2720472")

	//获取多个商品的价格
	goods2 := []string{"2720472", "1580032", "4861069"}
	getPriceOnly(goods2)
}
