package server

import (
	"log"
	"os"
	"path/filepath"
)

var (
	dataPath string
)

func initDataPath() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Getwd failed with err: %v", err)
	}
	dataPath = filepath.Join(dir, "data")
	err = os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		log.Fatalf("os.MkdirAll %s failed with err: %v", dataPath, err)
	}
}
