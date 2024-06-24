package utils

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUserIdFromJWT(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(nil, nil)

	userID := primitive.NewObjectID().Hex()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
	})

	// echo-jwt, gets it from cookies: Authorization= and set it to the user by default
	ctx.Set("user", token)

	result, err := GetUserIDFromJWT(ctx)

	assert.NoError(t, err)
	expectedId, _ := primitive.ObjectIDFromHex(userID)
	assert.Equal(t, expectedId, result)
}
