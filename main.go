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
		deleteCount   int64 = 0
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
	prompt := flag.Bool("prompt", true, "Disable prompt that displays before full execution.")

	flag.Parse()

	_, err := ValidateMageRoot(*mageRootPtr)
	if err != nil {
		color.Red(err.Error())
		return
	}

	_, err = ValidateDBCredentials(*userPtr, *dbNamePtr, *hostPtr)
	if err != nil {
		color.Red(err.Error())
		return
	}

	color.Yellow("Setting up database connection.")
	db, err = DbConnect(*userPtr, *passwordPtr, *hostPtr, *dbNamePtr)
	if err != nil {
		color.Red(err.Error())
		return
	}

	if *dummyData {
		err := GenerateDummyImageData(*mageRootPtr, *imageCount)
		if err != nil {
			color.Red(err.Error())
			return
		}
		color.Green("Dummy data has been generated successfully")
		return
	}

	if !*dryRunPtr && *prompt {
		result := FullExecutionPrompt(*dryRunPtr)
		if !result {
			color.Red("Aborting full execution")
			return
		}
	}

	color.Yellow("Collecting gallery records")
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

	color.Yellow("Collecting files to delete.")

	filesToDelete, totalFileSize, err := CollectFiles(files, *mageRootPtr, galleryValues, *includeCachePtr)
	if err != nil {
		color.Red(err.Error())
	}

	deleteMessage := DeleteMessage(*dryRunPtr)

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(6)

	color.Yellow("Start processing files.")
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

	g.Go(func() error {
		if *dryRunPtr {
			deleteCount, err = CountRecordsToDelete()
			if err != nil {
				color.Red(err.Error())
				return err
			}
		}

		return nil
	})

	g.Go(func() error {
		if *dryRunPtr {
			deleteCount, err = CountRecordsToDelete()
			if err != nil {
				color.Red(err.Error())
				return err
			}
		}

		return nil
	})

	g.Go(func() error {
		if !*dryRunPtr {
			deleteCount, err = DeleteGalleryRecords()
			if err != nil {
				color.Red(err.Error())
				return err
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	color.Green("Found " + strconv.Itoa(len(filesToDelete)) + " files for " + strconv.FormatFloat(totalFileSize/1024/1024, 'f', 2, 32) + " MB")
	color.Green("Found " + strconv.FormatInt(deleteCount, 10) + " database value(s) to remove")
}

func DeleteMessage(isDryRun bool) string {
	var deleteMessage string = "DRY-RUN: "

	if !isDryRun {
		deleteMessage = "REMOVING: "
	}

	deleteMessage = color.YellowString(deleteMessage)

	return deleteMessage
}
