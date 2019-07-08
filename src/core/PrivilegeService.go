package core

import "github.com/casbin/casbin"

type Menu struct {
	Title string
	Name  []string
}

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

func GetMenuByUser(user TUser) ([]LeftMenu, error) {
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
	err, pList := GetTreeLeftMenu()
	resultList := []LeftMenu{}
	for _, p := range pList {
		menu := LeftMenu{}
		menu.Name = p.Name
		if p.LeftMenu != nil {
			for _, a := range p.LeftMenu {
				for _, b := range list {
					if a.Redirecturl == b[1] {
						menu.LeftMenu = append(menu.LeftMenu, a)
						break
					}
				}
			}
		}
		if menu.LeftMenu != nil {
			resultList = append(resultList, menu)
		}
	}
	println(len(resultList))
	return resultList, nil
}

//递归查询页面权限
func (menu *Menu) FindShowHtml(p Privilege) {
	if len(p.Privilege) > 0 {
		for _, secondP := range p.Privilege {
			menu.FindShowHtml(secondP)
		}
	} else {
		if p.Type == 1 {
			menu.Name = append(menu.Name, p.Name)
		}
		return
	}
}
