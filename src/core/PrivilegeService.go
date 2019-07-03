package core

import "github.com/casbin/casbin"

//判断当前用户是否拥有path路径的权限
func PathExistPrivilege(loginUser string, path string) (bool, error) {
	confUrl := "./src/config/rbac_model.conf"
	csvUrl := "./src/config/basic_policy.csv"
	e := casbin.NewEnforcer(confUrl, csvUrl)
	sub := loginUser
	obj := path
	act := "read"
	if e.Enforce(sub, obj, act) == true {
		return true, nil
	} else {
		Logger(loginUser + "," + path + "," + "false")
		return false, nil
	}

}
