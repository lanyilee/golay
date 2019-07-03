package main

import (
	"controler"
	"log"
	"net/http"
)

func main() {
	//confUrl := "./src/config/rbac_model.conf"
	//csvUrl := "./src/config/basic_policy.csv"
	//e := casbin.NewEnforcer(confUrl, csvUrl)
	////sub := "alice"
	////obj := "data1"
	////act := "read"
	//if e.Enforce("alice","data1","read")==true{
	//	println(12)
	//	return
	//}
	//println(false)
	//return
	HandleRouter()
	err := http.ListenAndServe(":8888", nil)
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
