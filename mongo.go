package utility

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoClient struct
type MongoClient struct {
	IsConnected bool
	DB          *mongo.Database
}

var G_MONGO_CLIENT *MongoClient

//NewMongoClient func
func NewMongoClient() *MongoClient {
	mc := MongoClient{false, nil}
	Logger.Info("<< START CONNECT MONGODB >>")
	mongoURI := G_CONFIGER.GetString("mongo.uri")
	Logger.Info("mongo.uri :", mongoURI)
	// open connection
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		Logger.Error(err)
		return nil
	}
	// check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		Logger.Error(err)
		return nil
	}
	mongoDb := G_CONFIGER.GetString("mongo.db")
	Logger.Info("mongo.db :", mongoDb)
	mc.DB = client.Database(mongoDb)
	colList, _ := mc.DB.ListCollectionNames(context.TODO(), bson.D{{}})
	Logger.Info("mongo.collections : ", strings.Join(colList, ", "))
	mc.IsConnected = true
	Logger.Info("<< CONNECT SUCCESSFULLY >>")
	return &mc
}

//MongoFind func
func MongoFind(collection *mongo.Collection, sort interface{}, limit int64, filter interface{}) (results []bson.M) {
	opts := options.Find().SetSort(sort).SetLimit(limit)
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		Logger.Error(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		Logger.Error(err)
	}
	return
}

//MongoInsertMany func
func MongoInsertMany(collection *mongo.Collection, dataList []interface{}) {
	// set the Ordered option to false to allow both operations to happen even if one of them errors
	opts := options.InsertMany().SetOrdered(false)
	res, err := collection.InsertMany(context.TODO(), dataList, opts)
	if err != nil {
		Logger.Error(err)
	}
	fmt.Printf("inserted documents with IDs %v\n", res.InsertedIDs)
}

// MongoDistinct func
func MongoDistinct(collection *mongo.Collection, fieldName string, filter interface{}) []interface{} {
	opts := options.Distinct().SetMaxTime(0)
	values, err := collection.Distinct(context.TODO(), fieldName, filter, opts)
	if err != nil {
		Logger.Error(err)
	}
	return values
}

// MongoAggregate func
func MongoAggregate(collection *mongo.Collection, pipeline mongo.Pipeline) (results []bson.M) {
	opts := options.Aggregate().SetMaxTime(0)
	cursor, err := collection.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		Logger.Error(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		Logger.Error(err)
	}
	return
}
