package main

import (
	"controler"
	"core"
	"log"
	"net/http"
)

func main() {
	err, _ := core.GetXmlPrivileges()
	core.ErrHandler(err)
	HandleRouter()
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Panic(err)
	}

}

//路由转发
func HandleRouter() {

	http.Handle("/Golay.do", controler.Golay())
	http.Handle("/Login.do", controler.Login())
	http.Handle("/Privilege.do", controler.Privilege())

	http.Handle("/GetConfigPrivileges.do", controler.GetConfigPrivileges())
}
