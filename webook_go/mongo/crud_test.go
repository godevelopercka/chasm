package mongo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {
	// 控制初始化超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 每个命令（查询）执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		// 执行成功
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		// 执行失败
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	opts := options.Client().ApplyURI("mongodb://root:example@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	mdb := client.Database("webook")
	col := mdb.Collection("articles")
	defer func() {
		_, err = col.DeleteMany(ctx, bson.D{})
	}()

	res, err := col.InsertOne(ctx, Article{
		Id:      123,
		Title:   "我的标题",
		Content: "我的内容",
	})
	assert.NoError(t, err)
	// 这个是文档 ID,也就是 mongoDB 中的 _id 字段
	fmt.Printf("id %s", res.InsertedID)

	// bson
	// 找 ID = 123 的
	filter := bson.D{bson.E{Key: "id", Value: 123}}
	var art Article
	err = col.FindOne(ctx, filter).Decode(&art)
	assert.NoError(t, err)
	fmt.Printf("%#v \n", art)

	art = Article{}
	err = col.FindOne(ctx, Article{Id: 123}).Decode(&art)
	if err == mongo.ErrNoDocuments {
		fmt.Println("没有数据")
	}
	assert.NoError(t, err)
	fmt.Printf("%#v \n", art)

	sets := bson.D{bson.E{Key: "$set", Value: bson.E{Key: "title", Value: "新的标题"}}}
	updateRes, err := col.UpdateMany(ctx, filter, sets)
	assert.NoError(t, err)
	fmt.Println("affected", updateRes.ModifiedCount)
	updateRes, err = col.UpdateMany(ctx, filter, bson.D{
		bson.E{Key: "$set", Value: Article{Title: "我的标题2", AuthorId: 123456}}})
	assert.NoError(t, err)
	fmt.Println("affected", updateRes.ModifiedCount)

	// 写法一
	//or := bson.A{bson.D{bson.E{"id", 123}},
	//	bson.D{bson.E{"id", 456}}}
	// 写法二
	or := bson.A{bson.M{"id": 123}, bson.M{"id": 456}}
	orRes, err := col.Find(ctx, bson.D{bson.E{"$or", or}})
	assert.NoError(t, err)
	var ars []Article
	err = orRes.All(ctx, &ars)
	assert.NoError(t, err)

	and := bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"title", "我的标题2"}}}
	andRes, err := col.Find(ctx, bson.D{bson.E{"$and", and}})
	assert.NoError(t, err)
	ars = []Article{}
	err = andRes.All(ctx, &ars)
	assert.NoError(t, err)

	//in := bson.D{bson.E{"id", bson.D{bson.E{"$in", []any{123, 456}}}}}
	in := bson.D{bson.E{"id", bson.M{"in": []any{123, 456}}}}
	inRes, err := col.Find(ctx, in)
	assert.NoError(t, err)
	ars = []Article{}
	err = inRes.All(ctx, &ars)
	assert.NoError(t, err)

	inRes, err = col.Find(ctx, in, options.Find().SetProjection(bson.M{
		"id":    1,
		"title": 1,
	}))
	assert.NoError(t, err)
	ars = []Article{}
	err = inRes.All(ctx, &ars)
	assert.NoError(t, err)

	idxRex, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})
	assert.NoError(t, err)
	fmt.Println(idxRex)

	delRes, err := col.DeleteMany(ctx, filter)
	assert.NoError(t, err)
	fmt.Println("deleted", delRes.DeletedCount)

}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
