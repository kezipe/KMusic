package controllers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kezipe/musicserv/initializers"
	"github.com/kezipe/musicserv/models"
	"github.com/kezipe/musicserv/utils"
)

func SongsCreate(c *gin.Context) {
	// Get multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": "Multipart form required"})
		return
	}

	// Get title from form
	titles := form.Value["title"]
	if len(titles) == 0 {
		c.JSON(400, gin.H{"error": "Title is required"})
		return
	}
	title := titles[0]

	// Get the audio file
	files := form.File["audio"]
	if len(files) == 0 {
		c.JSON(400, gin.H{"error": "Audio file is required"})
		return
	}
	audioFile := files[0]

	// Open the file
	file, err := audioFile.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open audio file"})
		return
	}
	defer file.Close()

	// Generate a unique object name
	objectName := fmt.Sprintf("songs/%s-%s", uuid.New().String(), audioFile.Filename)

	// Upload to S3
	bucketName := os.Getenv("S3_BUCKET_NAME")
	err = utils.UploadFile(bucketName, objectName, file, audioFile.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Create a song record
	song := models.Song{
		Title:    title,
		AudioKey: objectName,
	}

	result := initializers.DB.Create(&song)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create song record: " + result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": song.Title + " created successfully",
		"song":    song,
	})
}

func SongsIndex(c *gin.Context) {
	var songs []models.Song
	result := initializers.DB.Find(&songs)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"songs": songs,
	})
}

func SongsShow(c *gin.Context) {
	var song models.Song
	id := c.Param("id")

	result := initializers.DB.First(&song, id)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"error": "Song not found",
		})
		return
	}

	// Generate pre-signed URL for the audio file
	bucketName := os.Getenv("S3_BUCKET_NAME")
	presignedURL, err := utils.GetPresignedURL(bucketName, song.AudioKey)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to generate audio URL: " + err.Error(),
		})
		return
	}

	// Create response with song data and presigned URL
	c.JSON(200, gin.H{
		"song": gin.H{
			"id":        song.ID,
			"title":     song.Title,
			"audioUrl":  presignedURL,
			"createdAt": song.CreatedAt,
			"updatedAt": song.UpdatedAt,
		},
	})
}

func SongsUpdate(c *gin.Context) {
	var body struct {
		Title    string `json:"title"`
		AudioKey string `json:"audio_key"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	id := c.Param("id")
	var song models.Song

	result := initializers.DB.First(&song, id)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"error": "Song not found",
		})
		return
	}

	song.Title = body.Title
	song.AudioKey = body.AudioKey

	initializers.DB.Model(&song).Updates(song)

	c.JSON(200, gin.H{
		"message": "Song updated successfully",
		"song":    song,
	})
}

func SongsDelete(c *gin.Context) {
	var song models.Song
	if err := initializers.DB.Where("id = ?", c.Param("id")).First(&song).Error; err != nil {
		c.JSON(404, gin.H{"error": "Song not found"})
		return
	}

	// Delete from S3 first
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if err := utils.DeleteObject(bucketName, song.AudioKey); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete audio file: " + err.Error()})
		return
	}

	// Then delete from database
	if err := initializers.DB.Delete(&song).Error; err != nil {
		// Note: at this point the S3 object is deleted but database deletion failed
		// In a production environment, you might want to implement a cleanup routine
		// to handle this kind of inconsistency
		c.JSON(500, gin.H{"error": "Failed to delete song record: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Song and associated audio file deleted successfully"})
}
