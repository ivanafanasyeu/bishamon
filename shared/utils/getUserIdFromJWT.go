package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserIDFromJWT(c echo.Context) (primitive.ObjectID, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return userObjectID, nil
}
