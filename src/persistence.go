package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 持久化接口,存对数据进行持久化
type Persistence interface {
	Save(documents []interface{}) (mongo.InsertManyResult, error)
}

// 创建mongo客户端
func CreateMonoPersistence() (Persistence, error) {
	mongoClient := new(mongoPersistence)
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	mongoClient.client = client
	return *mongoClient, err
}

type mongoPersistence struct {
	client *mongo.Client
}

func (mongo mongoPersistence) Save(documents []interface{}) (mongo.InsertManyResult, error) {
	collection := mongo.client.Database("douban").Collection("tv")
	res, err := collection.InsertMany(context.TODO(), documents)
	return *res, err
}
