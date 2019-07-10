package controler

import (
	"core"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Controler struct {
	config   core.Config
	redisCli *redis.Client
}

//Init Controler
func InitControler() *Controler {
	controler := Controler{}
	config, err := core.ReadConfig("./config.conf")
	if err != nil {
		core.Logger(err.Error())
	}
	redisCli, err := core.NewClientZero(config)
	if err != nil {
		core.Logger("InitControler:" + err.Error())
	}
	controler.config = config
	controler.redisCli = redisCli
	return &controler
}

//login
func Login() http.Handler {
	fn := func(response http.ResponseWriter, req *http.Request) {
		control := InitControler()
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
			if control.redisCli != nil {
				resb, err := core.CheckUserInRedis(control.redisCli, username, password)
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
						//存入token
						err = core.SaveToken(control.redisCli, token.String(), username)
						if err != nil {
							log.Panic(err)
						}
						resbBytes, err = json.Marshal(&resb)
						if err != nil {
							core.Logger(err.Error())
						}
					}
					response.Write(resbBytes)
					return
				}
			}
			db := core.CreatEngine(control.config)
			//数据库查询
			resb, err := core.Login(db, username, password, control.redisCli)
			if err != nil {
				log.Panic(err)
			}
			//生成token
			token, err := uuid.NewV4()
			if err != nil {
				log.Panic(err)
			}
			resb.Data = token.String()
			err = core.SaveToken(control.redisCli, token.String(), username)
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
	}
	return http.HandlerFunc(fn)
}

//privilege
func Privilege() http.Handler {
	fn := func(response http.ResponseWriter, req *http.Request) {
		control := InitControler()
		resp := core.ResponseBase{}
		//header
		token := req.Header.Get("GolayToken")
		if token == "" {
			resp.Message = "您尚未登录或登录已过期"
			resp.StatusCode = 402
		} else {
			loginUser, err := core.GetUserByToken(token, control.redisCli)
			if err != nil {
				core.Logger(err.Error())
			}
			if loginUser == nil {
				resp.Message = "您尚未登录或登录已过期"
				resp.StatusCode = 402
			} else {
				pathByte, err := ioutil.ReadAll(req.Body)
				if err != nil {
					core.Logger(err.Error())
				}
				path := string(pathByte[:])
				println(path)
				//通过user获取相应权限
				bo, err := core.PathExistPrivilege(loginUser.Username, path)
				if err != nil {
					core.Logger(err.Error())
				}
				if bo {
					resp.StatusCode = 200
					resp.Data = bo
				} else {
					resp.StatusCode = 403
					resp.Message = "你没有相关权限"
					resp.Data = bo
				}

			}
		}
		respBytes, err := json.Marshal(&resp)
		if err != nil {
			core.Logger(err.Error())
		}
		response.Write(respBytes)
	}
	return http.HandlerFunc(fn)
}

func Golay() http.Handler {
	fn := func(response http.ResponseWriter, req *http.Request) {

	}
	return http.HandlerFunc(fn)
}

func GetConfigPrivileges() http.Handler {
	fn := func(response http.ResponseWriter, req *http.Request) {
		resp := core.ResponseBase{}
		err, treeData := core.GetTreePrivileges()
		if err != nil {
			core.Logger(err.Error())
			resp.StatusCode = 500
		} else {
			resp.StatusCode = 200
			resp.Data = treeData
		}
		respBytes, err := json.Marshal(&resp)
		if err != nil {
			core.Logger(err.Error())
		}
		response.Write(respBytes)
	}
	return http.HandlerFunc(fn)
}

//获取菜单栏
func GetMenu() http.Handler {
	fn := func(response http.ResponseWriter, req *http.Request) {
		user, resp, err := GetUser(req)
		if err != nil {
			core.Logger(err.Error())
		} else if user.Username == "" {
			core.Logger("GetMenu:can't not found loginUser")
		} else {
			//通过用户获取用户相应菜单
			menuList, err := core.GetMenuByUser(user)
			if err != nil {
				resp.StatusCode = 500
				resp.Message = err.Error()
			} else {
				resp.StatusCode = 200
				resp.Data = menuList
				resp.Message = user.Realname
			}
		}
		respBytes, err := json.Marshal(&resp)
		if err != nil {
			core.Logger(err.Error())
		}
		response.Write(respBytes)

	}
	return http.HandlerFunc(fn)
}

//获取登录用户信息
func GetUser(req *http.Request) (core.TUser, core.ResponseBase, error) {
	control := InitControler()
	resp := core.ResponseBase{}
	user := core.TUser{}
	//header
	token := req.Header.Get("GolayToken")
	if token == "" {
		resp.Message = "您尚未登录或登录已过期"
		resp.StatusCode = 402
		return user, resp, nil
	} else {
		loginuser, err := core.GetUserByToken(token, control.redisCli)
		if err != nil {
			return user, resp, err
		}
		if loginuser == nil {
			resp.Message = "您尚未登录或登录已过期"
			resp.StatusCode = 402
		} else {
			user = *loginuser
		}
	}
	return user, resp, nil

}

//退出登录
func Logout() http.Handler {
	fn := func(response http.ResponseWriter, req *http.Request) {
		user, resp, err := GetUser(req)
		if err != nil {
			core.Logger(err.Error())
		} else if user.Username == "" {
			core.Logger("Logout:can't not found loginUser")
		} else {
			control := InitControler()
			db := core.CreatEngine(control.config)
			token := req.Header.Get("GolayToken")
			err = core.Logout(db, control.redisCli, &user, token)
			if err != nil {
				core.Logger(err.Error())
			}
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			core.Logger(err.Error())
		}
		core.Logger(user.Username + ":退出登录")
		response.Write(respBytes)
	}
	return http.HandlerFunc(fn)
}
