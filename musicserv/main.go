package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kezipe/musicserv/controllers"
	"github.com/kezipe/musicserv/initializers"
	"github.com/kezipe/musicserv/middleware"
	"github.com/kezipe/musicserv/models"
	"github.com/kezipe/musicserv/utils"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.DB.AutoMigrate(&models.Song{})

	if err := utils.InitS3(); err != nil {
		panic("Failed to initialize S3: " + err.Error())
	}
}

func main() {
	r := gin.Default()

	r.POST("/songs", middleware.RequireAuth, controllers.SongsCreate)
	r.GET("/songs", middleware.RequireAuth, controllers.SongsIndex)
	r.GET("/songs/:id", middleware.RequireAuth, controllers.SongsShow)
	r.PUT("/songs/:id", middleware.RequireAuth, controllers.SongsUpdate)
	r.DELETE("/songs/:id", middleware.RequireAuth, controllers.SongsDelete)
	r.GET("/amiauthorized", middleware.RequireAuth, controllers.AmIAuthorized)
	r.Run()
}
