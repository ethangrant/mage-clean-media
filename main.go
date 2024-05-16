package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"

	_ "github.com/go-sql-driver/mysql"
)

type File struct {
	Value string
	FullFilePath string
	FileSize int64
}

const mediaPath = "pub/media/catalog/product"

var db *sql.DB

func main() {
	var (
		files []File
		galleryValues []string
		filesToDelete []File
		totalFileSize float64 = 0
		deleteCount int
	)

	// handle arguments
	mageRootPtr := flag.String("mage-root", "", "Declare absolute path to the root of your magento installation")
	userPtr := flag.String("user", "", "Database username (required)")
	passwordPtr := flag.String("password", "", "Database password (required)")
	hostPtr := flag.String("host", "", "Database host (required)")
	dbNamePtr := flag.String("name", "", "Database name (required)")
	dryRunPtr := flag.Bool("dry-run", true, "Runs script without deleting files or DB records.")
	includeCachePtr := flag.Bool("no-cache", false, "Exclude files from catalog/product/cache directory.")

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

	// collect all files in pub/media/catalog/product
	absoluteMediaPath := *mageRootPtr + mediaPath

	filepath.WalkDir(absoluteMediaPath, func(path string, file fs.DirEntry, err error) error {
		var mediaFile File
		var fullPath string

		if err != nil {
			panic(err)
		}

		if !file.IsDir() {
			fileInfo, err := os.Stat(path)
			if err != nil {
				color.Red("There was a problem reading a file: " + err.Error())
			}
			fullPath = strings.Replace(path, absoluteMediaPath, "", -1)
			path = filepath.Base(path)
			mediaFile = File{
				Value: path,
				FullFilePath: fullPath,
				FileSize: fileInfo.Size(),
			}
            files = append(files, mediaFile)
        } 

		return nil;
	});

	// for _, file := range files {
	// 	fmt.Println("File: " + file.Value)
	// 	fmt.Println("File: " + file.FullFilePath)
	// }

	// collect all records from media tables
	const galleryValuesQuery = `
	SELECT gallery.value
	FROM catalog_product_entity_media_gallery AS gallery
	INNER JOIN catalog_product_entity_media_gallery_value_to_entity AS to_entity
	ON gallery.value_id = to_entity.value_id;`

	const placeholderQuery = `
	SELECT value FROM core_config_data WHERE path LIKE "%placeholder%" AND value IS NOT NULL;
	`

	const countDeleteQuery = `
	SELECT count(*) FROM catalog_product_entity_media_gallery AS gallery LEFT JOIN catalog_product_entity_media_gallery_value_to_entity AS to_entity ON gallery.value_id = to_entity.value_id WHERE (to_entity.value_id IS NULL);
	`

	const deleteGalleryQuery = `
	DELETE gallery FROM catalog_product_entity_media_gallery AS gallery
	LEFT JOIN catalog_product_entity_media_gallery_value_to_entity AS to_entity
	ON gallery.value_id = to_entity.value_id
	WHERE (to_entity.value_id IS NULL)
`

	rows, err := db.Query(galleryValuesQuery)
	if err != nil {
		color.Red("There was a problem collecting gallery records: " + err.Error())
		return
	}

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			color.Red("Error scanning row:", err)
			return
		}

		value = filepath.Base(value)

		galleryValues = append(galleryValues, value)
	}

	rows.Close()

	rows, err = db.Query(placeholderQuery)
	if err != nil {
		color.Red("There was a problem collecting placeholder image paths: " + err.Error())
		return
	}

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			color.Red("Error scanning row:", err)
			return
		}

		value = filepath.Base(value)

		galleryValues = append(galleryValues, value)
	}

	// for _, galleryValue := range galleryValues {
	// 	fmt.Println("Gallery value: " + galleryValue)
	// }

	// loop each file check if it appears in DB	
	for _, file := range files {
		var deleteFile bool = true

		if !*includeCachePtr {
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

	deleteMessage := DeleteMessage(*dryRunPtr)
	for _, file := range filesToDelete {
		if !*dryRunPtr {
			// Delete the files
		}

		fmt.Println(deleteMessage + file.FullFilePath)
	}

	// fmt.Println("Total files: " + strconv.Itoa(len(files)))


	color.Green("Found " + strconv.Itoa(len(filesToDelete)) + " files for " + strconv.FormatFloat(totalFileSize / 1024 / 1024, 'f', 2, 32) + " MB")

	rows, err = db.Query(countDeleteQuery)
	if err != nil {
		color.Red("There was a problem counting db records to be deleted: " + err.Error())
		return
	}

	for rows.Next() {
		err := rows.Scan(&deleteCount)
		if err != nil {
			color.Red("Error scanning row:", err)
			return
		}
	}

	// Run query to delete media records with no value
	if !*dryRunPtr {
		// _, err = db.Query(deleteGalleryQuery)
		// if err != nil {
		// 	color.Red("There was a problem removing DB records: " + err.Error())
		// 	return
		// }
	}

	color.Green("Found "+ strconv.Itoa(deleteCount) +" database value(s) to remove")
}