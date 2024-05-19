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
	g1, _ := errgroup.WithContext(ctx)
	g1.SetLimit(100)
	g2, _ := errgroup.WithContext(ctx)
	g2.SetLimit(5)

	color.Yellow(fmt.Sprintf("Generating %d images", count))

		g1.Go(func() error {
		for j := 0; j < count; j++ {

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

		}
		return nil

		})

	color.Yellow("Starting DB inserts")

	// @todo batch inserts
		g2.Go(func() error {
	for j := 0; j < count; j++ {

			filename, _ := RandomFileName(40, charset, seededRand)
			err := InsertGalleryRecord("/" + filename)
			if err != nil {
				color.Red("problem inserting dummy records: " + err.Error())
				return err
			}
		}

			return nil
		})


	if err := g1.Wait(); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	if err := g2.Wait(); err != nil {
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
