package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var supportedExtensions = map[string]bool{
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

const (
	concurrency = 10
)

func main() {
	folderPath := "."
	if err := validateFolder(folderPath); err != nil {
		log.Fatalf("Folder validation failed: %v", err)
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	var (
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, concurrency)
		counts    = struct {
			sync.Mutex
			success, failure int
		}{}
	)

	for _, file := range files {
		if !shouldProcess(file) {
			continue
		}

		wg.Add(1)
		go func(f os.DirEntry) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := renameFile(folderPath, f); err != nil {
				log.Printf("Rename failed: %v", err)
				counts.Lock()
				counts.failure++
				counts.Unlock()
				return
			}

			counts.Lock()
			counts.success++
			counts.Unlock()
		}(file)
	}

	wg.Wait()
	log.Printf("Complete. Success: %d, Failures: %d", counts.success, counts.failure)
}

func validateFolder(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path inaccessible: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}
	return nil
}

func shouldProcess(file os.DirEntry) bool {
	if file.IsDir() {
		return false
	}
	ext := strings.ToLower(filepath.Ext(file.Name()))
	return supportedExtensions[ext]
}

func renameFile(folderPath string, file os.DirEntry) error {
	uuid, err := generateUUID()
	if err != nil {
		return fmt.Errorf("uuid generation failed: %w", err)
	}

	oldPath := filepath.Join(folderPath, file.Name())
	ext := strings.ToLower(filepath.Ext(file.Name()))
	newPath := filepath.Join(folderPath, uuid+ext)

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("filesystem error: %w", err)
	}

	log.Printf("Renamed: %s -> %s", file.Name(), filepath.Base(newPath))
	return nil
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return hex.EncodeToString(b), nil
}