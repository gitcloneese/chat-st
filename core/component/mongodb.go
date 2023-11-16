package component

import (
	"context"
	"x-server/core/apollo"
	mg "x-server/core/pkg/database/mongdb"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDB() (db *mg.DB, err error) {
	var cfg struct {
		Client *mg.Config
	}

	err = paladin.Get(apollo.MongodbNS).UnmarshalTOML(&cfg)
	if err != nil {
		return
	}
	log.Debug("mongodb.txt %+v", cfg.Client)
	db, err = mg.NewMongoDB(cfg.Client)
	if err != nil {
		log.Error("NewMongoDB Error: %v", err)
		return
	}
	err = db.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return
	}

	return
}
