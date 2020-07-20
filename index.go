package omongo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//CreateIndexes 创建对应的索引,输入库和文档
//ensure是一个自己定义的格式比如 name_1_pwd_1 那么就是分配一个复合索引
//如果u_user_1 那么就是分配一个user的唯一索引，请注意，唯一索引我们默认只接受单键，如果有多键，程序执行不会成功
func (t *MongoDB) CreateIndexes(db, c string, ensure []string) error {
	//首先读取索引
	var index []bson.M
	ii := t.CDb(db, c).Indexes()
	cursor, err := ii.List(nil)
	if err != nil {
		return err
	}
	err = cursor.All(nil, &index)
	if err != nil {
		return err
	}
	//获得当前索引的map便于去重
	indexNow := map[string]bool{}
	for _, v := range index {
		nm := fmt.Sprint(v["name"])
		if nm == "_id_" {
			continue
		}
		indexNow[nm] = true
	}
	//获取需要写入索引的map方便去重
	indexNew := map[string]bool{}
	for _, v := range ensure {
		indexNew[v] = true
	}
	//这个函数
	f := func(s string) (bson.D, bool) {
		ss := strings.Split(s, "_")
		if len(ss) == 3 { //这样就是唯一索引
			return bson.D{{ss[1], s2i(ss[2])}}, true
		}
		if len(ss)%2 != 0 {
			return nil, false
		}

		rltD := bson.D{}
		for i := 0; i < len(ss); i += 2 {
			rltD = append(rltD, bson.E{Key: ss[i], Value: s2i(ss[i+1])})
		}
		return rltD, false
	}
	//这里增加当前没有的索引，并标记多余的当前索引
	for k := range indexNew {
		if _, ok := indexNow[k]; !ok {
			keys, uq := f(k)
			if keys == nil {
				continue
			}
			model := mongo.IndexModel{
				Keys:    keys,
				Options: options.Index().SetBackground(true).SetName(k).SetUnique(uq).SetSparse(uq),
			}
			_, err = ii.CreateOne(context.Background(), model)
			log.Println("AddEnsure:", k, err)
			if err != nil {
				return err
			}
		} else {
			indexNow[k] = false
		}
	}
	//这将删除多余的当前索引
	for k, v := range indexNow {
		if v {
			_, err := ii.DropOne(context.Background(), k)
			log.Println("DeleteEnsure:", k, err)
		}
	}

	return nil
}
