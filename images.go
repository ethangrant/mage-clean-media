package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
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

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(100)

	color.Yellow("Starting image creation")

	for j := 0; j < count; j++ {
		g.Go(func() error {
			filename, subDir := RandomFileName(20, charset, seededRand)
			fullpath := mediaPath + filename

			// Check dir exists before creating file
			if _, err := os.Stat(mediaPath + subDir); os.IsNotExist(err) {
				err = os.MkdirAll(mediaPath+subDir, os.ModePerm)
				if err != nil {
					color.Red(err.Error())
					return err
				}
			}

			source, err := os.Open("placeholder.jpg")
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

			color.Green(filename)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	color.Yellow("Starting DB inserts")

	ctx = context.Background()
	g, _ = errgroup.WithContext(ctx)
	g.SetLimit(5)

	// @todo batch inserts
	for j := 0; j < count; j++ {
		g.Go(func() error {
			filename, _ := RandomFileName(30, charset, seededRand)
			err := InsertGalleryRecord("/" + filename)
			if err != nil {
				color.Red("problem inserting dummy records: " + err.Error())
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
