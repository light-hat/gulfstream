package main

import (
	"context"
	"log"
	"os/exec"
	"path/filepath"
	// "time"

	"github.com/fsnotify/fsnotify"
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

	// Создание бакета если нужно
	bucketName := "hls"
	ctx := context.Background()
	if exists, _ := minioClient.BucketExists(ctx, bucketName); !exists {
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("Bucket creation failed: %v", err)
		}
	}

	// Запуск FFmpeg для трансляции
	cmd := exec.Command("ffmpeg",
		"-i", "rtmp://localhost:1935/live/{stream}",
		"-map", "0:v:0", "-c:v:0", "libx264", "-b:v:0", "2500k", "-preset", "veryfast",
		"-map", "0:a:0", "-c:a:0", "aac", "-b:a:0", "128k",
		"-f", "hls", "-hls_time", "4", "-hls_list_size", "10",
		"-hls_flags", "delete_segments+append_list",
		"-var_stream_map", "v:0,a:0",
		"-master_pl_name", "master.m3u8",
		"-hls_segment_filename", "/tmp/stream_%v_%03d.ts",
		"/tmp/stream_%v.m3u8",
	)

	if err := cmd.Start(); err != nil {
		log.Fatalf("FFmpeg start failed: %v", err)
	}

	// Мониторинг файлов и загрузка в MinIO
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	watcher.Add("/tmp")

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if filepath.Ext(event.Name) == ".ts" || filepath.Ext(event.Name) == ".m3u8" {
					_, objectName := filepath.Split(event.Name)
					_, err := minioClient.FPutObject(ctx, bucketName, objectName, event.Name, minio.PutObjectOptions{})
					if err != nil {
						log.Printf("Upload failed: %s - %v", event.Name, err)
					} else {
						log.Printf("Uploaded: %s", objectName)
					}
				}
			case err := <-watcher.Errors:
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	log.Println("Transcoder running")
	if err := cmd.Wait(); err != nil {
		log.Printf("FFmpeg exited: %v", err)
	}
}
