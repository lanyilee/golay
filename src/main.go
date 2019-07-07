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
	//list:=e.GetFilteredPolicy(0,"admin")
	//list3 := e.GetFilteredPolicy(2,"read")
	//list2 ,err:=core.GetStringListInterSection(list,list3)
	//println(len(list)+ len(list3)+len(list2))
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
	http.Handle("/GetMenu.do", controler.GetMenu())
	http.Handle("/Privilege.do", controler.Privilege())

	http.Handle("/GetConfigPrivileges.do", controler.GetConfigPrivileges())
}
