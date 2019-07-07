package core

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"
)

//[]string to []interface
func StringArrayToInterfaceArray(t []string) []interface{} {
	s := make([]interface{}, len(t))
	for i, v := range t {
		s[i] = v
	}
	return s
}

//求交集
func GetListInterSection(listA []interface{}, listB []interface{}) ([]interface{}, error) {
	err := &MyError{}
	var list []interface{}
	if len(listA) == 1 && len(listB) == 1 {
		err.CreateTime = time.Now().String()
		err.Message = "两个集合太小"
		return nil, err
	} else {
		for _, A := range listA {
			for _, B := range listB {
				byteA, err := json.Marshal(A)
				byteB, err := json.Marshal(B)
				if err != nil {
					return nil, err
				}
				mdA := md5.New()
				mdA.Write(byteA)
				aMd5 := hex.EncodeToString(mdA.Sum(nil))
				mdB := md5.New()
				mdB.Write(byteB)
				bMd5 := hex.EncodeToString(mdB.Sum(nil))
				if aMd5 == bMd5 {
					list = append(list, A)
					break
				}
			}
		}
	}
	return list, nil
}

//求[][]string交集
func GetStringListInterSection(listA [][]string, listB [][]string) ([][]string, error) {
	err := &MyError{}
	var list [][]string
	if len(listA) == 1 && len(listB) == 1 {
		err.CreateTime = time.Now().String()
		err.Message = "两个集合太小"
		return nil, err
	} else {
		for _, A := range listA {
			for _, B := range listB {
				byteA, err := json.Marshal(A)
				byteB, err := json.Marshal(B)
				if err != nil {
					return nil, err
				}
				mdA := md5.New()
				mdA.Write(byteA)
				aMd5 := hex.EncodeToString(mdA.Sum(nil))
				mdB := md5.New()
				mdB.Write(byteB)
				bMd5 := hex.EncodeToString(mdB.Sum(nil))
				if aMd5 == bMd5 {
					list = append(list, A)
					break
				}
			}
		}
	}
	return list, nil
}
