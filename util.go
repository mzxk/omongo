package omongo

import (
	"log"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//IsDuplicate 确认这个错误类型是不是唯一索引重复，如果是，返回true
//主要目的是为了用户名或者手机号码，或者email的唯一性，很好用
func IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "E11000 duplicate key")
}
func s2i(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		i = 0
	}
	return i
}

//ID 返回ID,如果输入的hex为空，那么返回一个新的iD
func ID(h string) primitive.ObjectID {
	if h == "" {
		return primitive.NewObjectID()
	}
	r, err := primitive.ObjectIDFromHex(h)
	if err != nil {
		log.Println(err)
	}
	return r
}
