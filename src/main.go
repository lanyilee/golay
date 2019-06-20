package main

import (
	"core"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	HandleRouter()
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Panic(err)
	}
}

//路由转发
func HandleRouter() {
	handler := http.HandlerFunc(Golay)
	http.Handle("/Golay.do", handler)
	handler2 := http.HandlerFunc(Login)
	http.Handle("/Login.do", handler2)
}

func Golay(response http.ResponseWriter, req *http.Request) {
	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	if dataByte != nil {
		a := string(dataByte[:])
		println(a)

	}
	return
}

//login
func Login(response http.ResponseWriter, req *http.Request) {
	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	var username, password string
	if dataByte != nil {
		reqStr := string(dataByte[:])
		reqs := strings.Split(reqStr, "&")
		for _, req := range reqs {
			field := strings.Split(req, "=")
			if field[0] == "username" {
				username = field[1]
			}
			if field[0] == "password" {
				password = field[1]
			}
		}
		config, err := core.ReadConfig("./config.conf")
		if err != nil {
			log.Panic(err)
		}
		//缓存查询
		redisCli, err := core.NewClientZero(config)
		if err != nil {
			core.Logger(err.Error())
		}
		if redisCli != nil {
			resb, err := core.CheckUserInRedis(redisCli, username, password)
			if err == nil && (resb.StatusCode == 200 || resb.StatusCode == 401) {
				resbBytes, err := json.Marshal(&resb)
				if err != nil {
					core.Logger(err.Error())
				}
				core.Logger("缓存登录")
				response.Write(resbBytes)
				return
			}
		}
		db := core.CreatEngine(config)
		//数据库查询
		resb, err := core.Login(db, username, password, redisCli)
		if err != nil {
			log.Panic(err)
		}
		resbBytes, err := json.Marshal(&resb)
		if err != nil {
			core.Logger(err.Error())
		}
		core.Logger("数据库登录")
		response.Write(resbBytes)
	}
	return
}
