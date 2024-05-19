package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Value        string
	FullFilePath string
	FileSize     int64
}

const mediaPath = "pub/media/catalog/product"

func CollectFiles(files []File, mageRootPath string, galleryValues []string, includeCache bool) ([]File, float64, error) {
	var totalFileSize float64
	absoluteMediaPath := mageRootPath + mediaPath

	err := filepath.WalkDir(absoluteMediaPath, func(path string, file fs.DirEntry, err error) error {
		var mediaFile File
		var fullPath string

		if err != nil {
			return err
		}

		if !file.IsDir() {
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			}
			fullPath = strings.Replace(path, absoluteMediaPath, "", -1)
			path = filepath.Base(path)
			mediaFile = File{
				Value:        path,
				FullFilePath: fullPath,
				FileSize:     fileInfo.Size(),
			}

			// Check should delete
			if ShouldDeleteFile(mediaFile, galleryValues, includeCache) {
				files = append(files, mediaFile)
				totalFileSize += float64(mediaFile.FileSize)
			}
		}

		return nil
	})

	if err != nil {
		return files, totalFileSize, err
	}

	return files, totalFileSize, nil
}

func ShouldDeleteFile(file File, galleryValues []string, includeCache bool) (result bool) {
	result = true

	if !includeCache {
		if strings.HasPrefix(file.FullFilePath, "/cache") {
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

func DeleteFile(mageRootPath string, filePath string) (err error) {
	err = os.Remove(mageRootPath + mediaPath + filePath)
	if err != nil {
		return err
	}

	return nil
}
