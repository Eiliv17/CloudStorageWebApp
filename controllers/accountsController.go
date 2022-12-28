package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Eiliv17/CloudStorageWebApp/initializers"
	"github.com/Eiliv17/CloudStorageWebApp/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("accounts")

	// request body struct
	body := struct {
		FullName string `form:"fullname" binding:"required"`
		BirthDay string `form:"birthday" binding:"required"`
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}{}

	// get email, birthday, fullname and password from req body
	err := c.ShouldBind(&body)
	if err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Missing some fields",
		})
		return
	}

	// check if the email already exist inside the database
	emailFilter := bson.D{primitive.E{Key: "email", Value: body.Email}}

	// check for email
	emailcount, _ := coll.CountDocuments(context.TODO(), emailFilter)
	if emailcount > 0 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Email already registered",
		})
		return
	}

	// hash the password with bcrypt
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// parse the birth day
	birthDay, err := time.Parse("2006-01-02", body.BirthDay)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// create the account
	timeNow := time.Now()

	userAccount := models.Account{
		ID:        primitive.NewObjectIDFromTimestamp(timeNow),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		FullName:  body.FullName,
		BirthDay:  birthDay,
		Email:     body.Email,
		Password:  string(hashedPass),
	}

	// stores the user inside the database
	_, err = coll.InsertOne(context.TODO(), userAccount)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
