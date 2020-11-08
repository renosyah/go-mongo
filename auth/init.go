package auth

import "go.mongodb.org/mongo-driver/mongo"

var (
	dbPool *mongo.Database
	option *Option
)

func Init(db *mongo.Database, opt *Option) {
	dbPool = db
	option = opt
}
