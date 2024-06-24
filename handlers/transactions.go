package handlers

import (
	"bishamon/db"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	AccountID     primitive.ObjectID  `json:"accountId" bson:"accountId"`
	Currency      string              `json:"currency" bson:"currency"`
	Date          primitive.DateTime  `json:"date" bson:"date"`
	Description   string              `json:"description" bson:"description"`
	Type          string              `json:"type" bson:"type"`
	Amount        float64             `json:"amount" bson:"amount"`
	DestinationID *primitive.ObjectID `json:"desctionationId,omitempty" bson:"destinationId,omitempty"`
}

func HandleCreateTransaction(c echo.Context) error {
	var transaction Transaction

	if err := c.Bind(&transaction); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if transaction.Type == "transfer" {
		if transaction.DestinationID == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing DestinationID")
		}
	}

	collection := db.Mongo.Database("bishamon").Collection("transactions")
	_, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transaction)
}
