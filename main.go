package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strconv"

	"github.com/fatih/color"

	_ "github.com/go-sql-driver/mysql"
)

// type File struct {
// 	Value string
// 	FullFilePath string
// 	FileSize int64
// }

// const mediaPath = "pub/media/catalog/product"

var db *sql.DB

func main() {
	var (
		files []File
		galleryValues []string
	)

	mageRootPtr := flag.String("mage-root", "", "Declare absolute path to the root of your magento installation")
	userPtr := flag.String("user", "", "Database username (required)")
	passwordPtr := flag.String("password", "", "Database password (required)")
	hostPtr := flag.String("host", "", "Database host (required)")
	dbNamePtr := flag.String("name", "", "Database name (required)")
	dryRunPtr := flag.Bool("dry-run", true, "Runs script without deleting files or DB records.")
	includeCachePtr := flag.Bool("no-cache", true, "Exclude files from catalog/product/cache directory.")

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

	if !*dryRunPtr {
		result := FullExecutionPrompt(*dryRunPtr)
		if !result {
			color.Red("Aborting full execution")
			return
		}
	}

	db, err = DbConnect(*userPtr, *passwordPtr, *hostPtr, *dbNamePtr)
	if err != nil {
		color.Red(err.Error())
		return
	}

	files = CollectFiles(files, *mageRootPtr)

	galleryValues, err = GalleryValues()
	if err != nil {
		color.Red(err.Error())
		return;
	}

	galleryValues, err = Placeholders(galleryValues)
	if err != nil {
		color.Red(err.Error())
		return;
	}

	filesToDelete, totalFileSize := FilesToDelete(files, galleryValues, *includeCachePtr)

	deleteMessage := DeleteMessage(*dryRunPtr)
	for _, file := range filesToDelete {
		if !*dryRunPtr {
			err = DeleteFile(*mageRootPtr, file.FullFilePath)
			if err != nil {
				color.Red(err.Error())
				return
			}
		}

		fmt.Println(deleteMessage + file.FullFilePath)
	}


	color.Green("Found " + strconv.Itoa(len(filesToDelete)) + " files for " + strconv.FormatFloat(totalFileSize / 1024 / 1024, 'f', 2, 32) + " MB")

	if !*dryRunPtr {
		err = DeleteGalleryRecords()
		if err != nil {
			color.Red(err.Error())
			return
		}
	}

	deleteCount, err := CountRecordsToDelete()
	if err != nil {
		color.Red(err.Error())
		return
	}

	color.Green("Found "+ strconv.Itoa(deleteCount) +" database value(s) to remove")
}