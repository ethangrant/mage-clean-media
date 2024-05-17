package main

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"math/rand"
	"os"
	"time"
)

func GenerateDummyImageData(mageRootPath string, count int) {
	var mediaPath string = mageRootPath + "pub/media/catalog/product/"

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	// source, err := os.Open("placeholder.jpg")
	// if err != nil {
	// 	color.Red(err.Error())
	// 	return
	// }

	for j := 0; j < count; j++ {
		filename, subDir := RandomFileName(20, charset, seededRand)
		fullpath := mediaPath + filename

		// Check dir exists before creating file
		if _, err := os.Stat(mediaPath + subDir); os.IsNotExist(err) {
			err = os.MkdirAll(mediaPath+subDir, os.ModePerm)
			if err != nil {
				color.Red(err.Error())
				return
			}
		}

		color.Green(fullpath)

		source, err := os.Open("placeholder.jpg")
		if err != nil {
			color.Red(err.Error())
			return
		}

		destination, err := os.Create(fullpath)
		if err != nil {
			color.Red("problem creating destination file: " + err.Error())
			return
		}

		color.Yellow(destination.Name())
		color.Yellow(source.Name())

		_, err = io.Copy(destination, source)
		if err != nil {
			color.Red(err.Error())
			return
		}

		source.Close()
		destination.Close()
	}

	// Generate dummy data to insert into 'catalog_product_entity_media_gallery'
}

func RandomFileName(length int, charset string, seededRand *rand.Rand) (string, string) {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	name := string(b) + ".jpg"

	firstChar := string(name[0])
	secondChar := string(name[1])

	subDir := fmt.Sprintf("%s/%s/", firstChar, secondChar)

	return subDir + name, subDir
}
