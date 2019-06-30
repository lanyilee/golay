package core

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
)

type PrivilegeResult struct {
	Privileges xml.Name    `xml:"privileges""`
	Id         string      `xml:"id""`
	Name       string      `xml:"name"`
	Privilege  []Privilege `xml:"privilege"`
}

type Privilege struct {
	Id        int         `xml:"id"`
	Name      string      `xml:"name"`
	Type      int         `xml:"type"`
	Selector  string      `xml:"selector"`
	Remarks   string      `xml:"remarks"`
	Privilege []Privilege `xml:"privilege"`
}

type PrivilegeTreeJson struct {
	Id       int    `json:"id"`
	Name     string `json:"title"`
	Selector string `json:"href"`
}

func GetXmlPrivileges() (error, string) {
	content, err := ioutil.ReadFile("./src/config/Privileges.xml")
	if err != nil {
		return err, ""
	}

	// xml 解析到result的结构中
	var result PrivilegeResult
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return err, ""
	}

	// sturct to json
	jsonString, err := ConvertResultToJsonString(result)
	Logger(jsonString)
	if err != nil {
		return err, ""
	}
	return nil, jsonString

}

func GetTreePrivileges() (error, []Privilege) {
	content, err := ioutil.ReadFile("./src/config/Privileges.xml")
	if err != nil {
		return err, nil
	}

	// xml 解析到result的结构中
	var result PrivilegeResult
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return err, nil
	}

	return nil, result.Privilege

}

func ErrHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func ConvertResultToJsonString(param interface{}) (string, error) {
	result, err := json.Marshal(param)
	if err != nil {
		Logger("返回结果struct转json出错：" + err.Error())
	}
	return string(result), err
}
