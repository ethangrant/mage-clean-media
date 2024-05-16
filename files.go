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

func CollectFiles(files []File, mageRootPath string) ([]File, error) {
	absoluteMediaPath := mageRootPath + mediaPath

	err := filepath.WalkDir(absoluteMediaPath, func(path string, file fs.DirEntry, err error) error {
		var mediaFile File
		var fullPath string

		if err != nil {
			panic(err)
		}

		if !file.IsDir() {
			fileInfo, err := os.Stat(path)
			if err != nil {
				panic(err)
			}
			fullPath = strings.Replace(path, absoluteMediaPath, "", -1)
			path = filepath.Base(path)
			mediaFile = File{
				Value:        path,
				FullFilePath: fullPath,
				FileSize:     fileInfo.Size(),
			}
			files = append(files, mediaFile)
		}

		return nil
	})

	if err != nil {
		return files, err
	}

	return files, nil
}

func FilesToDelete(files []File, galleryValues []string, includeCache bool) (filesToDelete []File, totalFileSize float64) {
	for _, file := range files {
		var deleteFile bool = true

		if !includeCache {
			if strings.HasPrefix(file.FullFilePath, "/cache") {
				continue
			}
		}

		for _, galleryValue := range galleryValues {
			if file.Value == galleryValue {
				deleteFile = false
				break
			}
		}

		if deleteFile {
			filesToDelete = append(filesToDelete, file)
			totalFileSize += float64(file.FileSize)
		}
	}

	return filesToDelete, totalFileSize
}

func DeleteFile(mageRootPath string, filePath string) (err error) {
	err = os.Remove(mageRootPath + mediaPath + filePath)
	if err != nil {
		return err
	}

	return nil
}
