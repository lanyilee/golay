package core

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"time"
)

//登录判断，权限授予
func Login(db *xorm.Engine, username string, password string, cli *redis.Client) (resp ResponseBase, err error) {
	tUser, err := GetUserByName(db, username)
	if err != nil {
		return resp, err
	}
	if tUser.Username == "" {
		resp.Message = "无此账号"
		resp.StatusCode = 400
	}
	md5Pass := MD532(password + tUser.Salt)
	if md5Pass != tUser.Password {
		resp.Message = "密码输入有误，请重新输入"
		Logger("错误密码：" + password)
		resp.StatusCode = 401
	} else {
		resp.Message = "登录成功"
		resp.StatusCode = 200
	}
	//入缓存
	userBytes, err := json.Marshal(&tUser)
	if err != nil {
		Logger(err.Error())
	}
	userStr := string(userBytes[:])
	_, err = cli.Set(username, userStr, time.Minute*60).Result()
	if err != nil {
		Logger(err.Error())
	}
	return resp, nil

}

func CheckUserInRedis(cli *redis.Client, username string, password string) (resp ResponseBase, err error) {
	userJson, err := cli.Get(username).Result()
	if err != nil {
		Logger(err.Error())
		return resp, err
	}
	userData := []byte(userJson)
	tUser := &TUser{}
	err = json.Unmarshal(userData, tUser)
	if err != nil {
		Logger(err.Error())
		return resp, err
	}
	if tUser.Username == "" {
		return resp, err
	}
	md5Pass := MD532(password + tUser.Salt)
	if md5Pass != tUser.Password {
		resp.Message = "密码输入有误，请重新输入"
		resp.StatusCode = 401
	} else {
		resp.Message = "登录成功"
		resp.StatusCode = 200
	}
	return resp, nil
}

func SaveToken(cli *redis.Client, token string, username string) error {
	_, err := cli.Set(token, username, time.Minute*60).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetUserByToken(token string, cli *redis.Client) (*TUser, error) {
	tUser := &TUser{}
	username, err := cli.Get(token).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			println(123)
		}
		return nil, err
	}
	if username != "" {
		userStr, err := cli.Get(username).Result()
		if err != nil {
			return nil, err
		}
		userData := []byte(userStr)

		err = json.Unmarshal(userData, tUser)
		if err != nil {
			return tUser, err
		}
	}
	return tUser, nil
}

func GetUserByName(db *xorm.Engine, username string) (TUser, error) {
	tUser := TUser{}
	_, err := db.Where("username='" + username + "'").Get(&tUser)
	if err != nil {
		return tUser, err
	}
	return tUser, nil
}
