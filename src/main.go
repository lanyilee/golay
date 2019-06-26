package main

import (
	"io/ioutil"
	"log"
	"net/http"
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
