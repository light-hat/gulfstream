package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// Подключение к MinIO
	minioClient, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("MinIO connection failed: %v", err)
	}

	r := gin.Default()

	r.GET("/:stream/:file", func(c *gin.Context) {
		stream := c.Param("stream")
		file := c.Param("file")
		objectName := stream + "/" + file

		// Генерация временной ссылки
		presignedURL, err := minioClient.PresignedGetObject(
			context.Background(),
			"hls",
			objectName,
			5*time.Minute,
			nil,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		c.Redirect(http.StatusFound, presignedURL.String())
	})

	log.Println("HLS Origin running on :8080")
	r.Run(":8080")
}
