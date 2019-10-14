package main

import (
	"fmt"

	"golangapi/db/elastics"
	"golangapi/db/mgo"

	"golangapi/routers"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"golangapi/middlewares"
)

/*
//CustomerHandler is struct
type CustomerHandler struct{}

//Initialize is cus init
func (h *CustomerHandler) Initialize() {

}
*/

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("port", "8083")

	e := echo.New()

	//Start MongoDB Connect
	//Hold Mongo lib, It slower than mgo lib client
	//mongo.ConnectMongo()

	//Start Mgo Connect
	mgo.ConnectMgo()

	//Start Elastics Connect
	elastics.ConnectES()

	// Start Router
	routers.Init(e)

	// Start Logger
	middlewares.InitLog()

	//h := &handler{}
	//h.Initialize(e)

	//h := &handler{}
	//e.POST("/login", h.login)

	port := fmt.Sprintf(":%v", viper.GetString("port"))
	e.Logger.Fatal(e.Start(port))
}
