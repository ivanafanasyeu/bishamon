package handlers

import (
	"bishamon/db"
	"bishamon/shared/types"
	"bishamon/shared/utils"
	"bishamon/views"
	"bishamon/views/components/forms"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleAddAccount(c echo.Context) error {
	ctx := c.Request().Context()
	var account types.Account

	if err := c.Bind(&account); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		templ := forms.FormAccountError("Bad request, check sending data")
		utils.RenderTempl(c, templ)
		return nil
	}

	userObjectID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		templ := forms.FormAccountError("User JWT could not be serialized")
		utils.RenderTempl(c, templ)
		return nil
	}

	if account.Currency == "" {
		account.Currency = "USD"
	}

	account.UserId = userObjectID
	collection := db.Mongo.Database("bishamon").Collection("accounts")
	res, err := collection.InsertOne(ctx, account)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		templ := forms.FormAccountError(err.Error())
		utils.RenderTempl(c, templ)
		return nil
	}

	account.ID = res.InsertedID.(primitive.ObjectID)

	c.Response().Header().Set("HX-Trigger", "accountCreated")
	c.Response().WriteHeader(http.StatusOK)
	return c.NoContent(http.StatusOK)
}

func HandleUpdateAccount(c echo.Context) error {
	ctx := c.Request().Context()
	var account types.AccountUpdate

	accountID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		templ := forms.FormAccountError("User JWT could not be serialized")
		utils.RenderTempl(c, templ)
		return nil
	}

	if err := c.Bind(&account); err != nil {
		log.Printf("Error binding account data: %v", err)
		c.Response().WriteHeader(http.StatusInternalServerError)
		errorTempl := forms.FormAccountError(err.Error())
		utils.RenderTempl(c, errorTempl)
		return nil
	}

	userObjectID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		templ := forms.FormAccountError("User JWT could not be serialized")
		utils.RenderTempl(c, templ)
		return nil
	}

	collection := db.Mongo.Database("bishamon").Collection("accounts")
	filter := bson.M{"_id": accountID, "userId": userObjectID}
	// @todo: refactor with custom UnmarshalJSON or something, due to unmarshall in Bind not skipping nil: https://github.com/golang/go/issues/22480
	updateAccount := bson.M{}
	if account.Name != nil && *account.Name != "" {
		updateAccount["name"] = account.Name
	}
	if account.Balance != nil {
		updateAccount["balance"] = account.Balance
	}
	if account.Currency != nil {
		updateAccount["currency"] = account.Currency
	}
	if account.IsInTotalBalance != nil {
		updateAccount["isInTotalBalance"] = account.IsInTotalBalance
	}
	update := bson.M{"$set": updateAccount}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Response().WriteHeader(http.StatusNotFound)
			templ := forms.FormAccountError("Account not found")
			utils.RenderTempl(c, templ)
			return nil
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		templ := forms.FormAccountError("Failed to update account")
		utils.RenderTempl(c, templ)
		return nil
	}

	if res.MatchedCount == 0 {
		c.Response().WriteHeader(http.StatusNotFound)
		templ := forms.FormAccountError("Account not found or not authorized to update")
		utils.RenderTempl(c, templ)
		return nil
	}

	c.Response().Header().Set("HX-Trigger", "accountUpdated")
	c.Response().WriteHeader(http.StatusOK)
	return c.NoContent(http.StatusOK)
}

func HandleDeleteAccount(c echo.Context) error {
	ctx := c.Request().Context()

	accountID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userObjectID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "JWT could not be serialized")
	}

	collection := db.Mongo.Database("bishamon").Collection("accounts")
	filter := bson.M{"_id": accountID, "userId": userObjectID}
	res, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, "Account not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete account")
	}

	if res.DeletedCount == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Account not found or not authorized to delete")
	}

	c.Response().WriteHeader(http.StatusOK)
	return c.NoContent(http.StatusOK)
}

func HandleGetAllAccounts(c echo.Context) error {
	ctx := c.Request().Context()

	userObjectID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "JWT could not be serialized")
	}

	collection := db.Mongo.Database("bishamon").Collection("accounts")
	filter := bson.M{"userId": userObjectID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve data from db")
	}
	defer cursor.Close(ctx)

	var accounts []types.Account
	if err = cursor.All(ctx, &accounts); err != nil {
		log.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode accounts")
	}

	listOnly := c.QueryParam("listOnly")

	if listOnly == "true" {
		return utils.RenderTempl(c, views.AccountsCards(accounts))
	}

	return utils.RenderTempl(c, views.Accounts(accounts))
}

func HandleGetAccountModal(c echo.Context) error {
	ctx := c.Request().Context()

	userObjectID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "JWT could not be serialized")
	}

	accountID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var account types.Account
	collection := db.Mongo.Database("bishamon").Collection("accounts")
	filter := bson.M{"_id": accountID, "userId": userObjectID}

	dbErr := collection.FindOne(ctx, filter).Decode(&account)
	if dbErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not account with provided id in database")
	}

	c.Response().Header().Set("HX-Trigger", "edit")
	return utils.RenderTempl(c, views.EditAccountDialog(&account))
}

func HandleGetAccountsTotalBalance(c echo.Context) error {
	ctx := c.Request().Context()

	userObjectID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "JWT could not be serialized")
	}

	collection := db.Mongo.Database("bishamon").Collection("accounts")
	filter := bson.M{"userId": userObjectID}
	pipeline := mongo.Pipeline{
		{{
			Key:   "$match",
			Value: filter,
		}},
		{{
			Key:   "$group",
			Value: bson.M{"_id": nil, "totalBalance": bson.M{"$sum": "$balance"}},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "MongoDB failed to aggregate")
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode balances")
	}

	if len(result) > 0 {
		totalBalance := result[0]["totalBalance"]
		return c.String(http.StatusOK, utils.FormatNumberToCurrency(totalBalance.(float64)))
	}

	return c.JSON(http.StatusOK, 0)
}
