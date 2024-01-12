package tools

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	pblogin "xy3-proto/login"
)

// 存储 account, player, token, platform
const (
	DB  = "test"
	COL = "user"
)

var (
	dbClient *mongo.Client
)

func DBClient() *mongo.Client {
	if dbClient == nil {
		mdb()
	}
	if dbClient == nil {
		fmt.Println("dbClient is nil")
	}
	return dbClient
}

// 暂时想把每次请求到的账户存储起来， 这个都不用每次都从新重建账户，
// 并且把当天的token存起来， 这样就不用每次都从新获取token了

func mdb() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Printf("mdb err:%v", err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Printf("mdb ping err:%v", err)
	}
	dbClient = client
	ensureIndex(client)
}

// User
// 做的简单点， 就存这些就可以了
type User struct {
	Account    string `bson:"account"`
	PlayerId   int64  `bson:"playerId"`
	Token      string `bson:"token"`
	CreateAt   int64  `bson:"createAt"`
	PlatformId int    `bson:"platformId"`
}

// 创建索引
func ensureIndex(db *mongo.Client) {
	database := db.Database(DB)
	collection := database.Collection(COL)
	indexModel := []mongo.IndexModel{

		{
			Keys: bson.D{
				{"account", 1},
			},
		},
		{
			Keys: bson.D{
				{"playerid", 1},
			},
		},
	}
	_, err := collection.Indexes().CreateMany(
		context.TODO(),
		indexModel)
	if err != nil {
		fmt.Printf("mdb ensureIndex err:%v", err)
	}
}

// SetDbPlayer
// player login token ttl两个小时
func SetDbPlayer(m map[string]*pblogin.LoginRsp) {
	if !useDB {
		return
	}
	if dbClient == nil {
		mdb()
	}
	if len(m) == 0 {
		return
	}
	mClient := DBClient()
	if mClient == nil {
		return
	}
	collection := mClient.Database(DB).Collection(COL)
	opt := options.Update()
	opt.SetUpsert(true)
	docs := make([]interface{}, 0, len(m))
	now := time.Now().Unix()
	for k, v := range m {
		docs = append(docs, &User{
			Account:    k,
			PlayerId:   v.PlayerID,
			Token:      v.PlayerToken,
			CreateAt:   now,
			PlatformId: PlatformId,
		})
	}
	_, err := collection.InsertMany(context.TODO(), docs, options.InsertMany())
	if err != nil {
		fmt.Printf("SetDbPlayer err: %v", err)
	}
}

// GetDBPlayer
// 获取是否有从db中获取玩家
func GetDBPlayer() bool {
	if !useDB {
		return false
	}
	if dbClient == nil {
		mdb()
	}
	mClient := DBClient()

	if mClient == nil {
		return false
	}
	collection := mClient.Database(DB).Collection(COL)
	var res []*User
	startTime := time.Now().Unix() - 3600*2
	cursor, err := collection.Find(context.TODO(), bson.D{
		{Key: "platformId", Value: PlatformId},
		{Key: "createAt", Value: bson.M{"$gt": startTime}},
	})
	if err != nil {
		fmt.Printf("GetDBPlayer find err: %v", err)
		return false
	}

	if err := cursor.All(context.TODO(), &res); err != nil {
		fmt.Printf("GetDBPlayer cursor All err: %v", err)
		return false
	}
	if len(res) == 0 {
		fmt.Printf("GetDBPlayer getPlayers empty")
		return false
	}

	for _, user := range res {
		GameLoginResp[user.Account] =
			&pblogin.LoginRsp{PlayerID: user.PlayerId, PlayerToken: user.Token}
	}

	fmt.Printf("GetDBPlayer getPlayers len:%v", len(GameLoginResp))
	return true
}
