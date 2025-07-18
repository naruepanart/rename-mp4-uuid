package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var exts = map[string]bool{
	// Image formats
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".webp": true,
	".tiff": true,
	".svg":  true,
	// Video formats
	".mp4":  true,
	".mov":  true,
	".avi":  true,
	".mkv":  true,
	".flv":  true,
	".wmv":  true,
	".webm": true,
	".mpeg": true,
	".mpg":  true,
	".3gp":  true,
}

func main() {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(f.Name()))
		if !exts[ext] {
			continue
		}

		uuid, err := generateUUID()
		if err != nil {
			log.Println(err)
			continue
		}

		newName := uuid + ext
		if err := os.Rename(f.Name(), newName); err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Renamed: %s -> %s", f.Name(), newName)
	}
}

func generateUUID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
