package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var (
		files         []File
		galleryValues []string
	)

	mageRootPtr := flag.String("mage-root", "", "Declare absolute path to the root of your magento installation")
	userPtr := flag.String("user", "", "Database username (required)")
	passwordPtr := flag.String("password", "", "Database password (required)")
	hostPtr := flag.String("host", "", "Database host (required)")
	dbNamePtr := flag.String("name", "", "Database name (required)")
	dryRunPtr := flag.Bool("dry-run", true, "Runs script without deleting files or DB records.")
	includeCachePtr := flag.Bool("no-cache", true, "Exclude files from catalog/product/cache directory.")
	dummyData := flag.Bool("dummy-data", false, "Set flag to generate a set of dummy image data.")
	imageCount := flag.Int("image-count", 500, "Define number of images to generate with dummy data option.")

	flag.Parse()

	_, err := ValidateMageRoot(*mageRootPtr)
	if err != nil {
		color.Red(err.Error())
		return
	}

	_, err = ValidateDBCredentials(*userPtr, *passwordPtr, *dbNamePtr, *hostPtr)
	if err != nil {
		color.Red(err.Error())
		return
	}

	db, err = DbConnect(*userPtr, *passwordPtr, *hostPtr, *dbNamePtr)
	if err != nil {
		color.Red(err.Error())
		return
	}

	if *dummyData {
		GenerateDummyImageData(*mageRootPtr, *imageCount)
		color.Green("Dummy data has been generated successfully")
		return
	}

	// if !*dryRunPtr {
	// 	result := FullExecutionPrompt(*dryRunPtr)
	// 	if !result {
	// 		color.Red("Aborting full execution")
	// 		return
	// 	}
	// }

	files, err = CollectFiles(files, *mageRootPtr)
	if err != nil {
		color.Red(err.Error())
	}

	galleryValues, err = GalleryValues()
	if err != nil {
		color.Red(err.Error())
		return
	}

	galleryValues, err = Placeholders(galleryValues)
	if err != nil {
		color.Red(err.Error())
		return
	}

	filesToDelete, totalFileSize := FilesToDelete(files, galleryValues, *includeCachePtr)

	deleteMessage := DeleteMessage(*dryRunPtr)

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	g.SetLimit(10000)

	for _, file := range filesToDelete {
		g.Go(func() error {
			if !*dryRunPtr {
				err = DeleteFile(*mageRootPtr, file.FullFilePath)
				if err != nil {
					return err
				}
			}
			fmt.Println(deleteMessage + file.FullFilePath)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	color.Green("Found " + strconv.Itoa(len(filesToDelete)) + " files for " + strconv.FormatFloat(totalFileSize/1024/1024, 'f', 2, 32) + " MB")

	// if !*dryRunPtr {
	// 	err = DeleteGalleryRecords()
	// 	if err != nil {
	// 		color.Red(err.Error())
	// 		return
	// 	}
	// }

	deleteCount, err := CountRecordsToDelete()
	if err != nil {
		color.Red(err.Error())
		return
	}

	color.Green("Found " + strconv.Itoa(deleteCount) + " database value(s) to remove")
}
