package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Eiliv17/CloudStorageWebApp/initializers"
	"github.com/Eiliv17/CloudStorageWebApp/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// signup controller
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

// login controller
func Login(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("accounts")

	// request body struct
	body := struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}{}

	// get email and password from req body
	err := c.ShouldBind(&body)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Missing some fields",
		})
		return
	}

	// look up requested user
	emailFilter := bson.D{primitive.E{Key: "email", Value: body.Email}}
	result := coll.FindOne(context.TODO(), emailFilter)

	// decode result
	var userAccount models.Account
	err = result.Decode(&userAccount)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Email or password wrong",
		})
		return
	}

	// compare sent pass with saved user pass hash
	err = bcrypt.CompareHashAndPassword([]byte(userAccount.Password), []byte(body.Password))
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Email or password wrong",
		})
		return
	}

	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userAccount.ID.Hex(),
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	// sign and get the complete encoded token as a string using the secret
	HMACSecret := os.Getenv("HMAC_SECRET")

	tokenString, err := token.SignedString([]byte(HMACSecret))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	// set token as cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	// redirects to dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}
