package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//获取签名
func GetCmsSign(config *Config, stringA string) (cmsSign string) {
	timeStemp := time.Now().Unix() * 1000
	timeStempStr := strconv.FormatInt(timeStemp, 10)
	nonstr := GetRandomString(16)
	token := "eyJhbGciOiJIUzI1NiIsInppcCI6IkRFRiJ9.eNp8zLsKwjAUgOF3OXMDuZKcbCWtULFavAyOaZtCHByaFATx3Y0v4Pz9_G945AgWUMzoZTBET1MgkjNPUBtGFF98GOmCQUmoIG1jidtXDuvTcsqQKibE0O7x0pjiMaXiru8ccXV3vx3JcKivu9O5_6HPYJlSWktpECvYUljj_Gf4-QIAAP__.qzZ-Ehwu6J2ZlXzmLQBQGQbidmc2NE8vE_aS_70E9n0"
	if len(stringA) != 0 {
		stringA += "&"
	}
	stringB := stringA + "key=" + MD532(timeStempStr+":"+nonstr+":"+token)
	sign := strings.ToUpper(MD532(stringB))

	cmsSignRaw := config.Account + "," + timeStempStr + "," + nonstr + "," + sign
	cmsSignBytes := []byte(cmsSignRaw)
	cmsSign = base64.StdEncoding.EncodeToString(cmsSignBytes)
	return cmsSign
}

//调用该接口获取人民日报往期数据
func GetPast(config *Config) (err error) {
	url := config.ApiAddress + "/past"
	cmsSign := GetCmsSign(config, "")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("CMS-SIGN", cmsSign)
	req.Header.Add("CMS-CLIENT", config.ClientID)
	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	rspData := &ResponseData{}
	rspBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		println(err)
	}
	if rspData.StatusCode == 200 {
		Logger(rspData.Message)
		pastDatas := []PastData{}
		json.Unmarshal(rspData.Data, &pastDatas)
		println(len(pastDatas))
	} else {
		Logger("出错！" + rspData.Message)
	}

	return nil
}

//报纸列表页和版面数据整合接口,参数pageSize为宽*高
func GetList(config *Config, date string, pagesSize string) (rmrbDatas []RmrbPage, err error) {
	url := config.ApiAddress + "/list?"
	param := ""
	if len(date) == 8 {
		param += "date=" + date
	}
	if len(pagesSize) == 0 {
		pagesSize = "400x571"
	}
	if param == "" {
		param += "pagesSize=" + pagesSize
	} else {
		param += "&pagesSize=" + pagesSize
	}
	url += param
	cmsSign := GetCmsSign(config, param)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("CMS-SIGN", cmsSign)
	req.Header.Add("CMS-CLIENT", config.ClientID)
	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	rspData := &ResponseData{}
	rspBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		println(err)
	}
	//rmrbDatas := [] RmrbPage{}
	if rspData.StatusCode == 200 {
		Logger(rspData.Message)
		json.Unmarshal(rspData.Data, &rmrbDatas)
		println(len(rmrbDatas))
		for i := 0; i < len(rmrbDatas); i++ {
			rmrbDataItems := []RmrbPageItem{}
			json.Unmarshal(rmrbDatas[i].Items, &rmrbDataItems)
			rmrbDatas[i].ItemList = &rmrbDataItems
			//println(rmrbDatas[i].ItemList)
		}

	} else {
		Logger("出错！" + rspData.Message)
	}

	return rmrbDatas, nil
}

//报纸列表页和版面数据整合接口,参数pageSize为宽*高
func GetDetail(config *Config, articleId string) (err error) {
	url := config.ApiAddress + "/detail?"
	param := "articleId=" + articleId
	url += param
	cmsSign := GetCmsSign(config, param)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("CMS-SIGN", cmsSign)
	req.Header.Add("CMS-CLIENT", config.ClientID)
	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	rspData := &ResponseData{}
	rspBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		println(err)
	}
	rmrbDetailDatas := RmrbNewsDetail{}
	if rspData.StatusCode == 200 {
		Logger(rspData.Message)
		json.Unmarshal(rspData.Data, &rmrbDetailDatas)
		println(rmrbDetailDatas.Content)
	} else {
		Logger("出错！" + rspData.Message)
	}
	return nil
}

//组合数据模板
func MakeUpModel(pointsStrs []string, urlList []string, picture string) (html string, err error) {
	file, err := ioutil.ReadFile("index.html")
	if err != nil {
		return "", err
	}
	fileReader := bytes.NewReader(file)
	doc, err := goquery.NewDocumentFromReader(fileReader)
	if err != nil {
		return "", err
	}
	doc.Find("img[usemap='#testmap']").RemoveAttr("src").SetAttr("src", picture)
	for i := 0; i < len(pointsStrs); i++ {
		shape := "rect"
		switch len(pointsStrs[i]) {
		case 3:
			shape = "circle"
			break
		case 4:
			shape = "rect"
			break
		default:
			shape = "poly"
		}
		area := "<area shape='" + shape + "' coords='" + pointsStrs[i] + "' href='" + urlList[i] + "' alt='" + strconv.Itoa(i) + "' target='_blank'/>\r\n"
		doc.Find("map[id='testmap']").AppendHtml(area)
	}
	html, err = doc.Html()
	if err != nil {
		return "", err
	}
	return html, nil
}
