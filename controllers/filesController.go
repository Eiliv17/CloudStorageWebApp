package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Eiliv17/CloudStorageWebApp/initializers"
	"github.com/Eiliv17/CloudStorageWebApp/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// file upload controller
func Upload(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("files")

	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"error": "Failed to upload file",
		})
		return
	}

	// get user info
	rawuser, exist := c.Get("user")
	if !exist {
		c.Redirect(http.StatusSeeOther, "/")
	}

	userAccount := rawuser.(models.Account)

	// save file in a folder
	filesavepath := "filedb/" + userAccount.ID.Hex() + "/" + file.Filename
	err = c.SaveUploadedFile(file, filesavepath)
	if err != nil {
		fmt.Println(filesavepath)
		fmt.Println(err)
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"error": "Failed to upload file",
		})
		return
	}

	// separate file name from file extension
	filextension := ""
	filename := file.Filename
	fileslice := strings.Split(file.Filename, ".")
	if len(fileslice) > 1 {
		filextension = fileslice[len(fileslice)-1]
		filename = strings.Join(fileslice[:len(fileslice)-1], "")
	}

	// create file model
	timeNow := time.Now()

	filedb := models.File{
		FileID:        primitive.NewObjectIDFromTimestamp(timeNow),
		User:          userAccount.ID,
		FileName:      filename,
		FileExtension: filextension,
		FileLocation:  "/" + filesavepath,
		InsertionDate: timeNow,
	}

	// insert file inside the db
	_, err = coll.InsertOne(context.TODO(), filedb)
	if err != nil {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

func Download(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("files")

	// get file id parameter
	fileID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// retrive file from database
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result := coll.FindOne(context.TODO(), filter)

	var file models.File
	err = result.Decode(&file)
	if err != nil {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"error": "File not found",
		})
		return
	}

	// get user info
	rawuser, exist := c.Get("user")
	if !exist {
		c.Redirect(http.StatusSeeOther, "/logout")
	}

	userAccount := rawuser.(models.Account)

	// check if file belongs to user
	if !(file.User == userAccount.ID) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"error": "File access not authorized",
		})
		return
	}

	c.FileAttachment("."+file.FileLocation, file.FileName+file.FileExtension)
}

func Dashboard(c *gin.Context) {

	c.HTML(http.StatusOK, "dashboard.html", gin.H{})
}
