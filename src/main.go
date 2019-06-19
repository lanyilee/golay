package main

import (
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
		println(username + password)

	}
	return
}
