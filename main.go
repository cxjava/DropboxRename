package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	log "github.com/cihub/seelog"
)

func init() {
	readConfig()
}

func main() {
	defer log.Flush()
	defer func() {
		if err := recover(); err != nil {
			log.Critical(err)
			log.Critical(string(debug.Stack()))
		}
	}()

	files := GetAllFiles(config.Dir)
	fmt.Println("===========================")
	for _, path := range files {
		newFileName := strings.Replace(path, config.FileName, "", -1)
		err := os.Rename(newFileName, newFileName+".dropbox.bak")
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = os.Rename(path, newFileName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(path + " to " + newFileName)
	}
}

func GetAllFiles(p string) (files []string) {
	walkDir, err := filepath.Abs(p)
	if err != nil {
		log.Error("GetAllFiles error:", err)
		return
	}
	log.Info(walkDir)
	files = make([]string, 0, 2000)
	err = filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error("Walk error:", err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, config.FileName) {
			files = append(files, path)
		}
		return nil
	})
	fmt.Println(len(files))
	return
}
