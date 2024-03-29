package core

import (
	"encoding/json"
	"time"
)

type TRmrbTemplate struct {
	Id         int    `xorm:"not null autoincr pk Int(11)"`
	Content    string `xorm:"Varchar(255)"`
	Updatetime string `xorm:"Varchar(25)"`
	Name       string `xorm:"Varchar(255)"`
}

//用户表
type TUser struct {
	Id            int       `xorm:"not null autoincr pk Int(11)"`
	Username      string    `xorm:"Varchar(255)"`
	Password      string    `xorm:"Varchar(255)"`
	Company       string    `xorm:"Varchar(255)"`
	Email         string    `xorm:"Varchar(255)"`
	Realname      string    `xorm:"Varchar(255)"`
	Lastlogintime time.Time `xorm:"datetime"`
	Salt          string    `xorm:"Varchar(255)"`
}

//返回
type ResponseBase struct {
	Data       interface{}
	Message    string
	StatusCode int
}

type RmrbPage struct {
	PageName string
	PageNum  string
	PagePic  string
	Items    json.RawMessage `json:"Items"`
	ItemList *[]RmrbPageItem
}

type RmrbPageItem struct {
	ArticleId     string
	CategoryId    string
	Id            string
	NewsDatetime  string
	NewsLink      string
	NewsTimestamp string
	PaperName     string
	PjCode        string
	Points        string
	RowNum        string
	SysCode       string
	Title         string
	ViewType      string
}

type RmrbNewsDetail struct {
	CategoryId    string
	Id            string
	SysCode       string
	ArticleId     string
	PaperName     string
	PjCode        string
	Title         string
	ShortTitle    string
	IntroTitle    string
	SubTitle      string
	Type          string
	Description   string
	Content       string
	Copyfrom      string
	Authors       string
	NewsTimestamp string
	NewsDatetime  string
	Medias        string
	Introduction  string
	NewsLink      string
	ShareUrl      string
	ShareSlogan   string
	Imgall        string
	Cover         string
}

type TPeUpdatestatus struct {
	Id     int `xorm:"not null pk int(5)"`
	Status int `xorm:"int(5)"`
}

type Phone struct {
	Mobile string
}
