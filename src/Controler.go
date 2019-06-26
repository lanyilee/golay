package main

import (
	"core"
	"encoding/json"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//login
func Login(response http.ResponseWriter, req *http.Request) {
	config, err := core.ReadConfig("./config.conf")
	if err != nil {
		log.Panic(err)
	}
	redisCli, err := core.NewClientZero(config)
	if err != nil {
		core.Logger(err.Error())
	}
	//header
	//reqToken:=req.Header.Get("golay_token")
	//core.GetUserByToken()
	//body
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
		//缓存查询
		if redisCli != nil {
			resb, err := core.CheckUserInRedis(redisCli, username, password)
			if err == nil && (resb.StatusCode == 200 || resb.StatusCode == 401) {
				resbBytes, err := json.Marshal(&resb)
				if err != nil {
					core.Logger(err.Error())
				}
				if resb.StatusCode == 200 {
					core.Logger("缓存登录成功")
					//生成token
					token, err := uuid.NewV4()
					if err != nil {
						log.Panic(err)
					}
					resb.Data = token.String()
					resbBytes, err = json.Marshal(&resb)
					if err != nil {
						core.Logger(err.Error())
					}
				}
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
		//生成token
		token, err := uuid.NewV4()
		if err != nil {
			log.Panic(err)
		}
		resb.Data = token.String()
		err = core.SaveToken(redisCli, token.String(), username)
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
