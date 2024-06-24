package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId           primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	Name             string             `json:"name" bson:"name"`
	Balance          float64            `json:"balance" bson:"balance"`
	Currency         string             `json:"currency" bson:"currency,omitempty"`
	IsInTotalBalance bool               `json:"isInTotalBalance" bson:"isInTotalBalance"`
}

type AccountUpdate struct {
	Name             *string  `json:"name,omitempty" bson:"name"`
	Balance          *float64 `json:"balance,omitempty" bson:"balance"`
	Currency         *string  `json:"currency,omitempty" bson:"currency"`
	IsInTotalBalance *bool    `json:"isInTotalBalance,omitempty" bson:"isInTotalBalance"`
}
