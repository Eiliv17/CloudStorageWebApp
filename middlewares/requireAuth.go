package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Eiliv17/CloudStorageWebApp/initializers"
	"github.com/Eiliv17/CloudStorageWebApp/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequireAuth(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("accounts")

	// get the cookie from request headers
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/accounts/login")
		return
	}

	// decode and validate it
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("HMAC_SECRET")), nil
	})

	if token == nil {
		c.Redirect(http.StatusSeeOther, "/accounts/login")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Redirect(http.StatusSeeOther, "/accounts/login")
			return
		}

		// find the user with token sub
		var userAccount models.Account
		returnedID, _ := primitive.ObjectIDFromHex(claims["userID"].(string))
		userIDFileter := bson.D{primitive.E{Key: "_id", Value: returnedID}}
		result := coll.FindOne(context.TODO(), userIDFileter)
		err := result.Decode(&userAccount)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/accounts/login")
			return
		}

		// attach user to the context
		c.Set("user", userAccount)

		// continue
		c.Next()

	} else {
		c.Redirect(http.StatusSeeOther, "/accounts/login")
		return
	}
}
