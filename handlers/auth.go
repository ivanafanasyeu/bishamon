package handlers

import (
	"bishamon/config"
	"bishamon/db"
	"bishamon/shared/utils"
	"bishamon/views"
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}

func getUser(username string) (*User, error) {
	collection := db.Mongo.Database("bishamon").Collection("users")

	var user User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func HandleUserLogin(c echo.Context) error {
	var body User

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := getUser(body.Username)
	if err != nil {
		c.Response().WriteHeader(http.StatusUnauthorized)
		templ := views.LoginError("Invalid credentials")
		utils.RenderTempl(c, templ)
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.Response().WriteHeader(http.StatusUnauthorized)
		templ := views.LoginError("Invalid credentials")
		utils.RenderTempl(c, templ)
		return nil
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time
	claims["sub"] = user.ID.Hex()

	conf := config.New()
	tokenString, err := token.SignedString([]byte(conf.App.JWT_TOKEN_SECRET))
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		templ := views.LoginError(err.Error())
		utils.RenderTempl(c, templ)
		return nil
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func HandlerUserLogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}
