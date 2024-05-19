package main

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

func DbConnect(user string, password string, host string, dbName string) (db *sql.DB, err error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	db, err = sql.Open("mysql", connection)
	if err != nil {
		return nil, errors.New("There was a problem connecting to the database: " + err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, errors.New("Could not ping database: " + err.Error())
	}

	return db, nil
}

func GalleryValues() (values []string, err error) {
	const sql = `
SELECT gallery.value
FROM catalog_product_entity_media_gallery AS gallery
INNER JOIN catalog_product_entity_media_gallery_value_to_entity AS to_entity
ON gallery.value_id = to_entity.value_id;`

	rows, err := db.Query(sql)
	if err != nil {
		return nil, errors.New("there was a problem collecting gallery records: " + err.Error())
	}

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		value = filepath.Base(value)

		values = append(values, value)
	}

	rows.Close()

	return values, nil
}

func Placeholders(galleryValues []string) (values []string, err error) {
	const sql = `
	SELECT value FROM core_config_data WHERE path LIKE "%placeholder%" AND value IS NOT NULL;
	`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, errors.New("there was a problem collecting placeholder records: " + err.Error())
	}

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		value = filepath.Base(value)

		galleryValues = append(galleryValues, value)
	}

	rows.Close()

	return galleryValues, nil
}

func CountRecordsToDelete() (count int64, err error) {
	const sql = `
	SELECT count(*) FROM catalog_product_entity_media_gallery AS gallery LEFT JOIN catalog_product_entity_media_gallery_value_to_entity AS to_entity ON gallery.value_id = to_entity.value_id WHERE (to_entity.value_id IS NULL);
	`
	rows, err := db.Query(sql)
	if err != nil {
		return 0, errors.New("There was a problem counting db records to be deleted " + err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, errors.New(err.Error())
		}
	}

	return count, nil
}

func DeleteGalleryRecords() (count int64, err error) {
	const sql = `
	DELETE gallery FROM catalog_product_entity_media_gallery AS gallery
	LEFT JOIN catalog_product_entity_media_gallery_value_to_entity AS to_entity
	ON gallery.value_id = to_entity.value_id
	WHERE (to_entity.value_id IS NULL)
	`
	result, err := db.Exec(sql)
	if err != nil {
		return 0, errors.New("There was a problem removing DB records: " + err.Error())
	}

	count, err = result.RowsAffected()
	if err != nil {
		return 0, errors.New(err.Error())
	}

	return count, nil
}

func InsertGalleryRecord(value string) error {
	const sql = `
	INSERT INTO catalog_product_entity_media_gallery (attribute_id, value, media_type, disabled) VALUES(88, ?, "image", 0)
	`

	_, err := db.Exec(sql, value)
	if err != nil {
		return err
	}

	return nil
}
