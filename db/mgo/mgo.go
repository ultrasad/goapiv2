package mgo

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

var (
	mongoDB *mgo.Session
	err     error
)

// ConnectMgo MongoDB Connect
func ConnectMgo() *mgo.Session {
	mongoHost := viper.GetString("mongo.host")
	mongoUser := viper.GetString("mongo.user")
	mongoPass := viper.GetString("mongo.pass")
	mongoDatabase := viper.GetString("mongo.db")

	connString := fmt.Sprintf("%v:%v@%v/%v", mongoUser, mongoPass, mongoHost, mongoDatabase)
	mongoDB, err = mgo.Dial(connString)
	if err != nil {
		log.Printf("dial mongodb server with connection string %q: %v", connString, err)
	}

	//BuildInfo
	//info, _ := mongoDB.BuildInfo()
	//fmt.Println("Mgo Info",info)

	return mongoDB
}

// MongoClient return MongoDB Session
func MongoClient() *mgo.Session {
	return mongoDB
}
