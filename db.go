package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Persistence 持久化接口,存对数据进行持久化
type Persistence interface {
	// Save 持久化多条结果
	Save(documents []interface{}) (*mongo.InsertManyResult, error)
	// SaveOne 持久化单条结果
	SaveOne(document interface{}) (*mongo.InsertOneResult, error)
}

// CreateMonoPersistence 创建mongo客户端
func CreateMonoPersistence() Persistence {
	mongoClient := new(mongoPersistence)
	ctx, cancer := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancer()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	mongoClient.collection = client.Database("douban").Collection("tv")
	return mongoClient
}

type mongoPersistence struct {
	collection *mongo.Collection
}

func (persistence *mongoPersistence) Save(documents []interface{}) (*mongo.InsertManyResult, error) {
	res, err := persistence.collection.InsertMany(context.TODO(), documents)
	return res, err
}

func (persistence *mongoPersistence) SaveOne(document interface{}) (*mongo.InsertOneResult, error) {
	res, err := persistence.collection.InsertOne(context.TODO(), document)
	return res, err
}
