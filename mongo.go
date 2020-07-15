package omongo

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoDB 1
type MongoDB struct {
	Clt *mongo.Client
	db  string
	url string
}

//NewMongoDB 初始化一个新的mongo实力
func NewMongoDB(url, db string) *MongoDB {
	r := &MongoDB{db: db}
	c, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Panic(err)
	}
	_ = c.Connect(nil)
	r.Clt = c

	return r
}

//C 获取默认的coll
func (t *MongoDB) C(c string) *Coll {
	return t.CDb(t.db, c)
}

//CDb 获取默认的coll
func (t *MongoDB) CDb(db, c string) *Coll {
	r := &Coll{}
	r.Collection = t.Clt.Database(db).Collection(c)
	return r
}

//Coll 1
type Coll struct {
	*mongo.Collection
}
