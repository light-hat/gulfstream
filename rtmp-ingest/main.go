package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/notedit/rtmp"
)

func main() {
	server := rtmp.NewServer(&rtmp.ServerConfig{
		Addr: ":1935",
	})

	server.HandlePublish = func(conn *rtmp.Conn) {
		streamKey := conn.URL.Path
		log.Printf("Stream started: %s", streamKey)
		
		// Пересылаем поток в трансодер
		transcoderURL := os.Getenv("TRANSCODER_URL") + streamKey
		relay := exec.Command("ffmpeg", "-i", "pipe:0", "-c", "copy", "-f", "flv", transcoderURL)
		
		relayIn, _ := relay.StdinPipe()
		relay.Stdout = os.Stdout
		relay.Stderr = os.Stderr
		
		go func() {
			if err := relay.Start(); err != nil {
				log.Printf("Relay failed: %v", err)
			}
		}()
		
		go func() {
			_, _ = io.Copy(relayIn, conn)
			relayIn.Close()
		}()
		
		if err := relay.Wait(); err != nil {
			log.Printf("Relay exited: %v", err)
		}
		
		log.Printf("Stream ended: %s", streamKey)
	}

	log.Println("RTMP Ingest running on :1935")
	log.Fatal(server.ListenAndServe())
}
