package main

import (
	"bishamon/config"
	"bishamon/db"
	"bishamon/handlers"
	"bishamon/shared/utils"
	"bishamon/views"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func indexHandler(c echo.Context) error {
	return utils.RenderTempl(c, views.Index())
}

func loginHandler(c echo.Context) error {
	return utils.RenderTempl(c, views.Login())
}

func main() {
	conf := config.New()

	errDb := db.InitMongo()
	if errDb != nil {
		fmt.Println(errDb) // log instead
		return
	}

	e := echo.New()
	e.Static("/", "assets")

	configJWT := echojwt.Config{
		SigningKey:  []byte(conf.App.JWT_TOKEN_SECRET),
		TokenLookup: "cookie:Authorization",
	}

	e.POST("/api/login", handlers.HandleUserLogin)
	e.POST("/api/logout", handlers.HandlerUserLogOut)
	e.GET("/login", loginHandler)

	e.GET("/", indexHandler, echojwt.WithConfig(configJWT))
	e.GET("/accounts", handlers.HandleGetAllAccounts, echojwt.WithConfig(configJWT))

	apiGroup := e.Group("/api", echojwt.WithConfig(configJWT))

	accountsGroup := apiGroup.Group("/accounts")
	accountsGroup.POST("", handlers.HandleAddAccount)
	accountsGroup.GET("/modal/:id", handlers.HandleGetAccountModal)
	accountsGroup.PATCH("/:id", handlers.HandleUpdateAccount)
	accountsGroup.DELETE("/:id", handlers.HandleDeleteAccount)
	accountsGroup.GET("/totalbalance", handlers.HandleGetAccountsTotalBalance)

	transactionsGroup := apiGroup.Group("/transactions")
	transactionsGroup.POST("", handlers.HandleCreateTransaction)

	e.Logger.Fatal(e.Start("127.0.0.1:" + conf.App.PORT))
}
