package core

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func CreatEngine(config Config) *xorm.Engine {
	dataSourceName := config.MysqlDataSource
	//整体格式:"数据库用户名:密码@(数据库地址:3306)/数据库实例名称?charset=utf8"
	//MysqlDataSource = root:root@(192.168.14.178:3306)/test?charset=gbk
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		Logger(err.Error())
	}
	//连接测试
	//if err := engine.Ping(); err!=nil{
	//	Logger(err.Error())
	//	return nil
	//}
	return engine
}

//插入数据
func InsertTRmrbTemplate(db *xorm.Engine, template TRmrbTemplate) error {
	_, err := db.Insert(template)
	return err
}

//查询数据是否存在
func IsExistTRmrbTemplate(db *xorm.Engine, template TRmrbTemplate) (bool, error) {
	temp := &TRmrbTemplate{}
	temp.Name = template.Name
	bool, err := db.Exist(temp)
	return bool, err
}

//查询是否pe新人数据(先查缓存，再查数据库)
func CheckIsNewUserByRedis(phone string, config Config) (int, error) {
	client, err := NewClientZero(config)
	if err != nil {
		return 0, err
	}
	re, err := client.Exists(phone).Result()
	if err != nil {
		return 0, err
	}
	bool := 0
	if int(re) == 0 {
		bool = 1
	}
	return bool, nil
}

//根据redisAB两库一起查询
func CheckIsNewUserByRedisAB(phone string, config Config) (int, error) {
	client, err := NewClientZero(config)
	if err != nil {
		return 0, err
	}
	re, err := client.Exists(phone).Result()
	if err != nil {
		return 0, err
	}
	clientB, err := NewClientOneByDefaultClient(client)
	if err != nil {
		return 0, err
	}
	reb, err := clientB.Exists(phone).Result()
	if err != nil {
		return 0, err
	}
	bool := 0
	if int(re) == 0 && int(reb) == 0 {
		bool = 1
	}
	return bool, nil
}

func CheckIsNewUserByRedisAB2(phone string, config Config) (int, error) {
	bool := 0
	client, errA := NewClientZero(config)
	if errA != nil {
		//0库崩了，查1库
		clientB, errB := NewClientOne(config)
		if errB != nil {
			return 0, errB
		}
		reb, err := clientB.Exists(phone).Result()
		if err != nil {
			return 0, err
		}
		if int(reb) == 0 {
			bool = 1
		}
		return bool, nil
	}
	reA, errA := client.Exists(phone).Result()
	if errA != nil {
		//0库崩了，查1库
		clientB, errB := NewClientOne(config)
		if errB != nil {
			return 0, errB
		}
		reb, err := clientB.Exists(phone).Result()
		if err != nil {
			return 0, err
		}
		if int(reb) == 0 {
			bool = 1
		}
		return bool, nil
	}
	//0库没崩，以0库为准
	if int(reA) == 0 {
		bool = 1
	}
	return bool, nil
}
