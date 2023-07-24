package routers

import (
	"net/http"
	"os"
	"time"

	"golangapi/controllers"
	"golangapi/middlewares"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//"net/http/httptrace"
)

// TimeEncoder return time encode
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02T15:04:05Z07:00"))
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z"))
}

//Init func
func Init(e *echo.Echo) {

	//hook := zapcore.AddSync(&middlewares.Logrus{Collection: "logger"})

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:   "time",
		LevelKey:  "level",
		NameKey:   "logger",
		CallerKey: "caller",
		//MessageKey:    "msg",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		//EncodeTime:    zapcore.ISO8601TimeEncoder,
		//EncodeTime:     zapcore.TimeEncoder(zapcore.PrimitiveArrayEncoder.AppendString(time.Time.Format("2006-01-02 15:04:05.000"))),
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		//zapcore.NewConsoleEncoder(NewEncoderConfig()),
		zapcore.NewJSONEncoder(encoderConfig),
		//zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), hook),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&middlewares.Logrus{Collection: "logger"})),
		//zapcore.NewMultiWriteSyncer(zapcore.AddSync(&middlewares.Logrus{Collection: "logger"})),
		//zap.DebugLevel,
		zap.InfoLevel,
	)

	zaplogger := zap.New(core, zap.AddCaller())

	e.Use(middlewares.ZapLogger(zaplogger))

	/*
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if username == "hanajung" && password == "secret" {
				return true, nil
			}
			return false, nil
		}))
	*/

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"level":"info", "time":"${time_rfc3339}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\r\n",
		//CustomTimeFormat: "2006-01-02T15:04:05Z",
		Output: &middlewares.Logs{Collection: "logs"},
		//Output: os.Stdout,
		//Output: echoLog,
	}))

	/*
		// Login route
		e.POST("/login", login)

		// Unauthenticated route
		e.GET("/", accessible)

		// Restricted group
		r := e.Group("/restricted")
		r.Use(middleware.JWT([]byte("secret")))
		r.GET("", restricted)
	*/

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/api")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("/restricted", restricted)
	r.GET("/todos", controllers.List)

	e.GET("/todos", controllers.List)
	e.POST("/todos", controllers.Create)
	e.GET("/todos/:id", controllers.View)
	e.PUT("/todos/:id", controllers.Done)
	e.DELETE("/todos/:id", controllers.Delete)

	e.GET("/allusers", controllers.GetAllUser)

	//GoRoutine
	e.GET("/hello", controllers.CallHelloRoutine)

	//Elastics Route
	//e.GET("/esversion", controllers.ESVersion)

	//Elastics Search
	e.GET("/essearch", controllers.Search)

	//e.Logger.Fatal(e.Start(port))

	//e.Get("/log", ...)
	//g := e.Group("/group", authenticationMiddleware)
	//g.Get("/auth", ...)
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "{status: Accessible}")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)

	//return c.String(http.StatusOK, "Welcome "+name+"!")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Welcome " + name,
	})
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	//email := c.FormValue("email")
	password := c.FormValue("password")

	// in our case, the only "valid user and password" is
	// user: rickety_cricket@example.com pw: shhh!
	// really, this would be connected to any database and
	// retrieving the user and validating the password
	//if email != "rickety_cricket@example.com" || password != "shhh!" {
	if username != "hanajung" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	// add any key value fields to the token
	//claims["email"] = "rickety_cricket@example.com"
	claims["name"] = "Hanajung"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	// return the token for the consumer to grab and save
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
