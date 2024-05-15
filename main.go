package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type File struct {
	Value string
	FullFilePath string
}

var db *sql.DB

const mediaPath = "pub/media/catalog/product"

func main() {
	var files []File
	var galleryValues []string
	var filesToDelete []File
	// handle arguments
	mageRoot := os.Args[1]
	isDryrun, err := strconv.ParseBool(os.Args[2])
	includeCache, err := strconv.ParseBool(os.Args[3])
	user := os.Args[4]
	password := os.Args[5]
	host := os.Args[6]
	dbName := os.Args[7]
	fmt.Println(mageRoot)
	fmt.Println(isDryrun)
	fmt.Println("connection details:")
	fmt.Println(user)
	fmt.Println(password)
	fmt.Println(host)
	fmt.Println(dbName)

	// read db creds from env.php


	// setup mysql connection
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}

	fmt.Println("Connected!")

	// collect all files in pub/media/catalog/product
	absoluteMediaPath := mageRoot + mediaPath

	fmt.Println("media path: " + absoluteMediaPath)

	filepath.WalkDir(absoluteMediaPath, func(path string, file fs.DirEntry, err error) error {
		var mediaFile File
		var fullPath string

		if err != nil {
			return err
		}

		if !file.IsDir() {
			fullPath = strings.Replace(path, absoluteMediaPath, "", -1)
			path = filepath.Base(path)
			mediaFile = File{
				Value: path,
				FullFilePath: fullPath,
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

	rows, err := db.Query(galleryValuesQuery)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			fmt.Println("Error scanning row:", err)
		}

		value = filepath.Base(value)

		galleryValues = append(galleryValues, value)
	}

	rows.Close()

	rows, err = db.Query(placeholderQuery)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			fmt.Println("Error scanning row:", err)
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
		}
	}

	for _, file := range filesToDelete {
		fmt.Println("Delete file: " + file.FullFilePath)
	}

	fmt.Println("Total files: " + strconv.Itoa(len(files)))
	fmt.Println("Total gallery values: " + strconv.Itoa(len(galleryValues)))
	fmt.Println("Total files deleted: " + strconv.Itoa(len(filesToDelete)))
	// NOT in db we remove the file
	// Run query to delete media records with no value
}