package core

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	//配置文件要通过tag来指定配置文件中的名称
	//api
	ApiAddress string `ini:"ApiAddress"`
	Account    string `ini:"Account"`
	ClientID   string `ini:"ClientID"`
	//MySql
	MysqlDataSource string `ini:"MysqlDataSource"`
	//
	FixedTime string `ini:"FixedTime"`
	//Redis
	RedisIP  string `ini:"RedisIP"`
	RedisPwd string `ini:"RedisPwd"`
}

func Logger(strContent string) {
	logPath := "./log/" + time.Now().Format("2006-01-02") + ".txt"
	file, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	fileTime := time.Now().Format("2006-01-02 15:04:05")
	fileContent := strings.Join([]string{"===", fileTime, "===", strContent, "\n"}, "")
	buf := []byte(fileContent)
	file.Write(buf)
	defer file.Close()
}

func ReadLog(path string, rex string, rexls string) string {
	file, _ := os.Open(path)
	defer file.Close()
	bye, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}
	str := string(bye[:])
	index := strings.Index(str, rex)
	index2 := strings.Index(str, rexls)
	if index > 0 {
		selected := str[index:index2]
		fmt.Println(selected)
		return selected
	}
	return ""
}

//读取配置文件并转成结构体
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path) //加载配置文件
	if err != nil {
		Logger("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config) //解析成结构体
	if err != nil {
		Logger("mapto config file fail!")
		return config, err
	}
	return config, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//获取目录
func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

//截取字符串,截取的不包括第end位
func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		panic("start is wrong")
	}
	if end < start || end > length {
		panic("end is wrong")
	}
	return string(rs[start:end])
}

//生成32位MD5
func MD532(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

//用linux自带的openssl加密3DES-CBC,command的首参是openssl,不是平常的/bin/bash
func Encrypt3DESByOpenssl(key string, fileName string) (desPath string, err error) {
	filePath := "./formatFiles/" + fileName
	desPath = filePath + ".des"
	fmt.Println("将要加密的文件地址：" + filePath)
	cmd := exec.Command("openssl", "enc", "-des-ede3-cbc", "-e", "-k", key, "-in", filePath, "-out", desPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Logger("Error:can not obtain stdout pipe for command")
		return "", err
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		Logger("Error:The command is err")
		return "", err
	}
	//读取所有输出
	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		Logger("wait error")
		return "", err
	}
	Logger("encrypt success:")
	//fmt.Printf("stdout:\n\n %s", "")
	return desPath, nil
}

//重构文本-模板
//func FormatJKText(kd *KdcheckResult) string {
//	accounttype := ""
//	if kd.accounttype == "1" {
//		accounttype = "手机宽带"
//	} else if kd.accounttype == "2" {
//		accounttype = "裸宽宽带"
//	} else {
//		accounttype = "其他宽带"
//	}
//	isYearPackAge := "否"
//	if kd.IsYearPackAge == "1" {
//		isYearPackAge = "是"
//	}
//	//str := "START|" + kd.KdAccount + "|\n" + "宽带属性|" + accounttype + "~家庭宽带~" + kd.UserStatus + "~" + isYearPackAge + "~" + kd.LastDate + "~" + kd.BroadSpeed + "|010000\nEND\n"
//	//脱敏
//	//phoneNum:=subString(kd.KdAccount,0,3)+"****"+subString(kd.KdAccount,7,11)
//	str := "START|" + kd.KdAccount + "|\n" + "宽带属性|" + kd.KdAccount + "~" + accounttype + "~" + kd.UserStatus + "~" + isYearPackAge + "~" + kd.LastDate + "~" + kd.BroadSpeed + "|010000\nEND\n"
//	//utf8->gbk
//	str = Encode(str)
//	return str
//}

//tar命令解压
func UnCompressFile(formatFilePath string) error {
	Logger("uncompressFile Path:" + formatFilePath)
	cmd := exec.Command("tar", "-xzvf", formatFilePath, "-C", "./files")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Logger("UnCompressFile Error:can not obtain stdout pipe for command")
		return err
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		Logger("UnCompressFile Error:The command is err")
		return err
	}
	//读取所有输出
	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		Logger("UnCompressFile wait error")
		fmt.Println(err)
		return err
	}
	Logger("UnCompressFile success:" + formatFilePath)
	return nil
}

func TarFile(formatFilePath string) error {
	Logger("uncompressFile Path:" + formatFilePath)
	cmd := exec.Command("sh", "./tar.sh", formatFilePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Logger("UnCompressFile Error:can not obtain stdout pipe for command")
		return err
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		Logger("UnCompressFile Error:The command is err")
		return err
	}
	//读取所有输出
	str, err := ioutil.ReadAll(stdout)
	mes := string(str)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if strings.Contains(mes, "success") {
			Logger("UnCompressFile ftp file success:" + formatFilePath)
			return nil
		}
		Logger("wait error:" + string(str) + ";")
		return err
	}
	Logger("UnCompressFile success:" + formatFilePath)
	return nil
}

//删除文件
func RemoveFiles(path string) error {
	Logger("RemoveFile Path:" + path)
	cmd := exec.Command("sh", "./rm.sh", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Logger("RemoveFiles Error:can not obtain stdout pipe for command")
		return err
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		Logger("RemoveFiles Error:The command is err")
		return err
	}
	//读取所有输出
	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		Logger("RemoveFiles wait error")
		fmt.Println(err)
		return err
	}
	Logger("RemoveFile success:" + path)
	return nil
}

//取定时时间
func GetFixTime(config *Config) (fixTime time.Time, err error) {
	fixTimeStr := config.FixedTime
	//fixTime := time.Date(2018, 11, 06, 07, 52, 0, 0, time.Local)
	year, err := strconv.Atoi(subString(fixTimeStr, 0, 4))
	if err != nil {
		return fixTime, err
	}
	monthNum, _ := strconv.Atoi(subString(fixTimeStr, 4, 6))
	if err != nil {
		return fixTime, err
	}
	day, _ := strconv.Atoi(subString(fixTimeStr, 6, 8))
	if err != nil {
		return fixTime, err
	}
	hour, _ := strconv.Atoi(subString(fixTimeStr, 8, 10))
	if err != nil {
		return fixTime, err
	}
	min, _ := strconv.Atoi(subString(fixTimeStr, 10, 12))
	if err != nil {
		return fixTime, err
	}
	//这个month竟然还是个time.Month类型，奇葩
	month := time.Month(monthNum)
	fixTime = time.Date(year, month, day, hour, min, 0, 0, time.Local)
	return fixTime, nil
}

//压缩文件
func CompressFile(formatFilePath string) error {
	//zipPath ="./formatFiles/"+zipPath
	cmd := exec.Command("gzip", formatFilePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Logger("compressFile Error:can not obtain stdout pipe for command")
		return err
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		Logger("compressFile Error:The command is err")
		return err
	}
	//读取所有输出
	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		Logger("compressFile wait error")
		fmt.Println(err)
		return err
	}
	Logger("compressFile success:" + formatFilePath + ".gz")
	return nil
}

func SyncLoggerNum(strContent string) {
	go func(str string) {
		defer func() {
			recover()
		}()
		Logger(str)
	}(strContent)
}

//取随机字符串
func GetRandomString(num int) string {
	chars := [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u", "v", "w", "s", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	if num < 0 {
		return ""
	}
	randStr := ""
	for i := 0; i < num; i++ {
		index := GenerateRangeNum(0, 35)
		randStr += chars[index]
	}
	return randStr
}

//睡眠一毫秒为：防止时间戳一样，取出的随机数是一样的情况出现
func GenerateRangeNum(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := r.Intn(max-min) + min
	time.Sleep(time.Millisecond * 1)
	return randNum
}
