package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(Golay)
	http.Handle("/Golay", handler)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Panic(err)
	}
}

//路由转发
func Golay(response http.ResponseWriter, req *http.Request) {
	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	if dataByte != nil {
		println(len(dataByte))

	}
	return
}
