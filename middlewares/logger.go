package middlewares

import (
	"encoding/json"
	"fmt"

	//"context"
	//"golangapi/db/mongo"
	"golangapi/db/mgo"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	collection string
	err        error
	fLogger    *lumberjack.Logger
	//errLog     *log.Logger
)

type (
	//Logger struct logger from go log
	Logger struct {
		Time       time.Time `bson:"time" json:"time"`
		Lv         string    `bson:"level" json:"level"`
		Prefix     string    `bson:"prefix" json:"prefix"`
		Message    string    `bson:"-" json:"message"`
		Data       CtxLogger `bson:"data" json:"data"`
		Collection string    `bson:"-"`
	}

	//CtxLogger struct logger req,res
	CtxLogger struct {
		ID  string      `bson:"id" json:"id" `
		Req interface{} `bson:"request" json:"req"`
		Res interface{} `bson:"response" json:"res"`
	}

	// Logs struct log from echo
	Logs struct {
		//ID           string    `json:"id" bson:"id"`
		//ID           primitive.ObjectID    `json:"id" bson:"_id"`
		Time         time.Time `bson:"time" json:"time"`
		RemoteIP     string    `bson:"remote_ip" json:"remote_ip"`
		Host         string    `bson:"host" json:"host"`
		Method       string    `bson:"method" json:"method"`
		URI          string    `bson:"uri" json:"uri"`
		Status       int       `bson:"status" json:"status"`
		Latency      int       `bson:"latency" json:"latency"`
		LatencyHuman string    `bson:"latency_human" json:"latency_human"`
		BytesIn      int       `bson:"bytes_in" json:"bytes_in"`
		BytesOut     int       `bson:"bytes_out" json:"bytes_out"`
		Collection   string    `bson:"-"`
	}

	//Logrus struct log from Logrus
	Logrus struct {
		// Time       time.Time   `json:"time" bson:"time"`
		// Animal     string      `json:"animal" bson:"animal"`
		// Data       ctxLogger   `bson:"data" json:"data"`
		// ID         string      `json:"id" bson:"id"`
		// Req        interface{} `json:"req" bson:"request"`
		// Res        interface{} `json:"res" bson:"response"`
		// Message    string      `bson:"-" json:"message"`
		// Collection string      `bson:"-"`
		Time       time.Time `bson:"time" json:"time"`
		Lv         string    `bson:"level" json:"level"`
		Prefix     string    `bson:"prefix" json:"prefix"`
		Message    string    `bson:"-" json:"message"`
		Data       CtxLogger `bson:"data" json:"data"`
		Collection string    `bson:"-"`
	}
)

//Init log
func init() {

	//some time shutdown database, you will need this.
	year, month, day := time.Now().Date()
	/*fLogger = &lumberjack.Logger{
		Filename: filepath.Join("./logs", strconv.Itoa(year)+"-"+strconv.Itoa(int(month))+"-"+strconv.Itoa(day)+".log"),
		MaxSize:  650,  // megabytes
		MaxAge:   15,   //days
		Compress: true, // disabled by default
	}*/

	fLogger = &lumberjack.Logger{
		Filename: filepath.Join("./logs", strconv.Itoa(year)+"-"+strconv.Itoa(int(month))+"-"+strconv.Itoa(day)+".log"),
		MaxSize:  650,  // megabytes
		MaxAge:   15,   //days
		Compress: true, // disabled by default
	}

	fmt.Println("init logs...")
}

func (lg *Logger) Write(logByte []byte) (n int, err error) {

	//fmt.Println("log", lg)
	//fmt.Println("logByte", logByte)

	err = json.Unmarshal(logByte, &lg)
	if err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %q", lg)
		return len(logByte), nil
		//return err
	}
	/*
		if err != nil {
			fmt.Println("\n err Logger, json Unmarshal >>>", err)
			return
		}
	*/

	//fmt.Println("lg.Message => ", lg.Message)
	//fmt.Println("lg.Data => ", lg.Data)

	err = json.NewDecoder(strings.NewReader(lg.Message)).Decode(&lg.Data)

	if err != nil {
		//fmt.Println("\n err json decode >>>", err)
		return
	}

	/* MongoClient */
	/*
		go func() {
			client := mongo.ClientManager()
			// create a new context with a 10 second timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			insertResult, err := client.Database("document").Collection(lg.Collection).InsertOne(ctx, &lg)
			if err != nil {
				fmt.Printf("\n err time:%s,%s\n", time.Now(), lg.Message)
			}
			fmt.Println("Inserted a Logger: ", insertResult.InsertedID)
		}()
	*/

	//MgoClient
	go func() {
		client := mgo.MongoClient().Copy()
		defer client.Close()

		if err := client.DB("document").C(lg.Collection).Insert(&lg); err != nil {
			fmt.Printf("\n err time:%s,%s\n", time.Now(), lg.Message)
		} else {
			//fmt.Printf("\n not err, time:%s\n", time.Now())
		}
	}()

	return len(logByte), nil
}

/*func LoggerLumberjack() *lumberjack.Logger {
	return fLogger
}*/

//echo Logs
// 2019-05-14, comment fix test
func (lg *Logs) Write(logEcho []byte) (n int, err error) {

	err = json.Unmarshal(logEcho, &lg)
	if err != nil {
		fmt.Println("\n err Logs, json Unmarshal >>>", err)
		return
	}

	//fLogger
	go func() {
		fLogger.Write(logEcho)
	}()

	//fmt.Printf("\n &lg Logs: %#v\n", &lg)
	// MongoClinet
	/*
		go func() {
			client := mongo.ClientManager()
			// create a new context with a 10 second timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			insertResult, err := client.Database("document").Collection(lg.Collection).InsertOne(ctx, &lg)
			if err != nil {
				fmt.Printf("\n err time:%s,%s\n", time.Now(), err)
			}
			fmt.Println("Inserted a Log: ", insertResult.InsertedID)
		}()
	*/

	//MgoClient
	go func() {
		client := mgo.MongoClient().Copy()
		defer client.Close()

		if err := client.DB("document").C(lg.Collection).Insert(&lg); err != nil {
			fmt.Printf("\n err Logs time:%s, %s\n", time.Now(), err)
		} else {
			//fmt.Printf("\n not err, time:%s\n", time.Now())
		}
	}()

	return len(logEcho), nil
}

//Logrus Log
func (lg *Logrus) Write(logByte []byte) (n int, err error) {

	//fmt.Println("logByte => ", logByte)
	//fmt.Println("&lg => ", &lg)
	//f := map[string]interface{}{}
	//var f interface{}
	//err = json.Unmarshal(logByte, &f)
	//m := f.(map[string]interface{})
	//foomap := m["foo"]
	//fmt.Println("F => ", &f)

	//err = json.Unmarshal(logByte, &f)
	err = json.Unmarshal(logByte, &lg)
	if err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}

		log.Printf("response: %q", lg)
		return len(logByte), nil
		//return err
	}

	//fmt.Println("Logrus lg => ", lg)
	//fmt.Println("lg Interface Data => ", &lg)
	//fmt.Println("lg.Message => ", lg.Message)
	//fmt.Println("lg.Data => ", lg.Data)

	//err = json.NewDecoder(strings.NewReader(lg.Message)).Decode(&lg.Data)
	err = json.NewDecoder(strings.NewReader(lg.Message)).Decode(&lg.Data)

	if err != nil {
		fmt.Println("\n err json decode >>>", err)
		return
	}

	//fmt.Printf("Message: %s Data: %s", lg.Message, lg.Data)
	//fmt.Printf("Response: %s", lg.Data)
	//fmt.Printf("lg.Collection: %s", lg.Collection)

	//fLogger
	go func() {
		fLogger.Write(logByte)
	}()

	//MgoClient
	go func() {
		client := mgo.MongoClient().Copy()
		defer client.Close()

		if err := client.DB("document").C(lg.Collection).Insert(&lg); err != nil {
			fmt.Printf("\n err Logs time:%s, %s\n", time.Now(), err)
		} else {
			//fmt.Printf("\n not err, time:%s\n", time.Now())
		}
	}()

	return len(logByte), nil
}
