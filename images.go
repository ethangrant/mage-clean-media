package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
)

func GenerateDummyImageData(mageRootPath string, count int) {
	var mediaPath string = mageRootPath + "pub/media/catalog/product/"

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(100)

	color.Yellow(fmt.Sprintf("Generating %d images", count))

	for j := 0; j < count; j++ {
		g.Go(func() error {
			filename, subDir := RandomFileName(40, charset, seededRand)
			fullpath := mediaPath + filename

			// Check dir exists before creating file
			if _, err := os.Stat(mediaPath + subDir); os.IsNotExist(err) {
				err = os.MkdirAll(mediaPath+subDir, os.ModePerm)
				if err != nil {
					color.Red(err.Error())
					return err
				}
			}

			source, err := os.Open("images/placeholder1.jpg")
			if err != nil {
				color.Red(err.Error())
				return err
			}

			destination, err := os.Create(fullpath)
			if err != nil {
				color.Red("problem creating destination file: " + err.Error())
				return err
			}

			_, err = io.Copy(destination, source)
			if err != nil {
				color.Red(err.Error())
				return err
			}

			source.Close()
			destination.Close()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	color.Yellow("Starting DB inserts")

	var fileNames []string
	for j := 0; j < count; j++ {
			filename, _ := RandomFileName(40, charset, seededRand)
			fileNames = append(fileNames, filename)
	}

	ctx = context.Background()
	g, _ = errgroup.WithContext(ctx)
	g.SetLimit(5)

	chunks := ChunkSlice(fileNames, 1000)
	for _, chunk := range chunks {
		g.Go(func() error {
			err := InsertMultipleGalleryRecords(chunk)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
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
