package main

import (
	"net/http"

	"github.com/Eiliv17/CloudStorageWebApp/controllers"
	"github.com/Eiliv17/CloudStorageWebApp/initializers"
	"github.com/Eiliv17/CloudStorageWebApp/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.LoadDatabase()
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("views/*")
	r.Static("/public", "./public")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// accounts routes
	racc := r.Group("/accounts")
	{
		// serves register html
		racc.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", gin.H{})
		})

		// handles the post request for registering an account
		racc.POST("/register", controllers.Signup)

		// serves login html
		racc.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{})
		})

		// handles the post request for logging in an account
		racc.POST("/login", controllers.Login)
	}

	r.GET("/dashboard", middlewares.RequireAuth, func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{})
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", "", 0, "", "", false, true)
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
