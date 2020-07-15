package omongo

import (
	"context"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
这个文件重构了一些原生驱动没有的易用性函数
*/
//BatchResult 用于之前mongodb驱动的ALL
type BatchResult struct {
	cursor *mongo.Cursor
	ctx    context.Context
	err    error
}

//All 用于mgo的all函数,这里使用了反射,输入必须是一个&[]struct
func (t *BatchResult) All(result interface{}) error {
	if t.err != nil {
		return t.err
	}
	resultV := reflect.ValueOf(result) //获取result对应反射的值,这通常是一个指针
	slicev := resultV.Elem()           //获取result指针的值,这应该是一个slice
	if slicev.Kind() != reflect.Slice {
		return errors.New("mustSliceType")
	}
	i := 0
	slicev = slicev.Slice(0, slicev.Cap()) //这为slice的值确认了一个空间,但是通常这个cap都是0...好像毫无意义这一句
	elemt := slicev.Type().Elem()          //slice的type为interface的[],然后在elem为result的struct
	defer t.cursor.Close(t.ctx)
	for {
		temp := reflect.New(elemt) //通过刚才获取到到result到结构注册一个新变量并取得指针
		if !t.cursor.Next(t.ctx) {
			break
		}
		e := t.cursor.Decode(temp.Interface()) //获取temp的interface以便于写入
		if e != nil {
			return e
		}
		slicev = reflect.Append(slicev, temp.Elem()) //slicev是一个值的[],temp是一个指针,所以要获取elem,在使用反射包的增加数组
		i++
	}
	if i == 0 {
		return nil
	}
	resultV.Elem().Set(slicev.Slice(0, i))
	return nil
}

//FindAll 1
func (t *Coll) FindAll(ctx context.Context, filter interface{}, opts ...*options.FindOptions) *BatchResult {
	if filter == nil {
		filter = bson.M{}
	}
	cr, err := t.Find(ctx, filter, opts...)
	br := &BatchResult{cr, ctx, err}
	return br
}

//Upsert 封装的updateone
func (t *Coll) Upsert(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return t.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
}
