package core

import "github.com/casbin/casbin"

func NewCasbinEnforcer() *casbin.Enforcer {
	confUrl := "./src/config/rbac_model.conf"
	csvUrl := "./src/config/basic_policy.csv"
	e := casbin.NewEnforcer(confUrl, csvUrl)
	return e
}

//判断当前用户是否拥有path路径的权限
func PathExistPrivilege(loginUser string, path string) (bool, error) {
	e := NewCasbinEnforcer()
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

func GetMenuByUser(user TUser) ([]string, error) {
	//获取user所有权限
	e := NewCasbinEnforcer()
	//list:=e.GetPermissionsForUser(user.Username)
	list1 := e.GetFilteredPolicy(0, "admin")
	list2 := e.GetFilteredPolicy(2, "read")
	list, err := GetStringListInterSection(list1, list2)
	if err != nil {
		Logger(err.Error())
		return nil, err
	}
	err, pList := GetTreePrivileges()
	resultList := [][]string{}
	for _, p := range pList {
		for _, mp := range list {
			if p.Selector == mp[1] {
				resultList = append(resultList, p.Selector)
			}
		}
	}
	println(len(resultList))
	return resultList, nil
}
