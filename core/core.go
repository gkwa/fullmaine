package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/go-logr/logr"
)

type FileProcessor interface {
	ProcessFiles(dir string, startNum int, logger logr.Logger) error
}

type fileProcessor struct{}

func NewFileProcessor() FileProcessor {
	return &fileProcessor{}
}

func (fp *fileProcessor) ProcessFiles(dir string, startNum int, logger logr.Logger) error {
	if err := ensureDir(dir, logger); err != nil {
		return err
	}

	re := regexp.MustCompile(`test_(\d{5})\.(md|golden)`)
	highestNum := make(map[string]int)
	filesExist := false

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			filesExist = true
			matches := re.FindStringSubmatch(info.Name())
			if len(matches) == 3 {
				num, _ := strconv.Atoi(matches[1])
				fileType := matches[2]
				if num > highestNum[fileType] {
					highestNum[fileType] = num
				}
			}
		}

		return nil
	})
	if err != nil {
		logger.Error(err, "Error walking the path", "dir", dir)
		return err
	}

	if !filesExist || startNum > highestNum["md"] {
		return createFiles(dir, startNum, logger)
	}

	newNum := highestNum["md"] + 100
	return createFiles(dir, newNum, logger)
}

func ensureDir(dir string, logger logr.Logger) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Info("Creating directory", "dir", dir)
		return os.MkdirAll(dir, 0o755)
	}
	return nil
}

func createFiles(dir string, num int, logger logr.Logger) error {
	for _, fileType := range []string{"md", "golden"} {
		if err := createFile(dir, num, fileType, logger); err != nil {
			return err
		}
	}
	return nil
}

func createFile(dir string, num int, fileType string, logger logr.Logger) error {
	fileName := fmt.Sprintf("test_%05d.%s", num, fileType)
	filePath := filepath.Join(dir, fileName)

	_, err := os.Create(filePath)
	if err != nil {
		logger.Error(err, "Error creating file", "path", filePath)
		return err
	}
	logger.Info("Created file", "path", filePath)
	return nil
}
