package main

import (
	"fmt"
	"github.com/MichaelTJones/walk"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Value        string
	FullFilePath string
	FileSize     int64
}

func DeleteFiles(files []File, mageRootPath string, galleryValues []string, includeCache bool, isDryRun bool) (int64, float64, error) {
	const mediaPath = "pub/media/catalog/product"
	var totalFileSize float64
	var deletedCount int64
	absoluteMediaPath := mageRootPath + mediaPath
	deleteMessage := DeleteMessage(isDryRun)

	err := walk.Walk(absoluteMediaPath, func(root string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// Get just the file name to compare against DB records
			path := filepath.Base(root)

			mediaFile := File{
				Value:        path,
				FullFilePath: root,
				FileSize:     info.Size(),
			}

			if ShouldDeleteFile(mediaFile, galleryValues, includeCache) {

				if !isDryRun {
					err = DeleteFile(root)
					if err != nil {
						return err
					}
				}

				fmt.Println(deleteMessage + root)
				totalFileSize += float64(mediaFile.FileSize)
				deletedCount++
			}
		}

		return nil
	})

	if err != nil {
		return deletedCount, totalFileSize, err
	}

	return deletedCount, totalFileSize, nil
}

func ShouldDeleteFile(file File, galleryValues []string, includeCache bool) (result bool) {
	result = true
	
	if !includeCache {
		if strings.Contains(file.FullFilePath, "catalog/product/cache") {
			// File is in cache dir, don't delete
			return false
		}
	}

	for _, galleryValue := range galleryValues {
		if file.Value == galleryValue {
			// File is in DB don't delete
			return false
		}
	}

	return result
}

func DeleteFile(path string) (err error) {
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func DeleteMessage(isDryRun bool) string {
	var deleteMessage string = "DRY-RUN: "

	if !isDryRun {
		deleteMessage = "REMOVING: "
	}

	deleteMessage = color.YellowString(deleteMessage)

	return deleteMessage
}
