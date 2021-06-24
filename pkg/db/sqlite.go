package db

import (
	"fmt"
	"gora/pkg/utils"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ImageDB struct {
	db *sqlx.DB
}

func NewImageDB() *ImageDB {
	db, err := NewConnSQLite()
	if err != nil {
		log.Println(err)
	}
	return &ImageDB{db}
}

// NewConnSQLite создаём подключение к нашей базе данных
func NewConnSQLite() (*sqlx.DB, error) {
	return sqlx.Connect("sqlite3", "ex.db")
}

type SQLiteDB struct {
	CRUDimages
}

type CRUDimages interface {
	GetImagePreview(id string) (string, error)
	GetImages() (*Images, error)
	AddImage(string) (bool, error)
	DeleteImage(id string) (bool, string, string, error)
}

func NewSQLiteDB() *SQLiteDB {
	return &SQLiteDB{
		CRUDimages: NewImageDB(),
	}
}

type Image struct {
	ID         string `json:"id"`
	PeviewFile string `json:"-"`
	Preview    string `json:"pre"`
}

type Images struct {
	Images []Image `json:"images"`
}

// GetImages возврящает имя файла
func (i *ImageDB) GetImagePreview(id string) (string, error) {
	var (
		path string
	)
	err := i.db.Get(&path, "select preview_name from images where id=$1", id)
	if err != nil {
		return path, err
	}
	return path, nil
}

// GetImages возврящает id и имя файлов-миниатюр
func (i *ImageDB) GetImages() (*Images, error) {
	imgs := []Image{}
	rows, err := i.db.Query("select id, preview_name from images")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		img := Image{}
		err := rows.Scan(&img.ID, &img.PeviewFile)
		if err != nil {
			log.Panicln(err)
			continue
		}
		fmt.Println(img)
		imgs = append(imgs, img)
	}
	return &Images{imgs}, nil
}

// DeleteImage возвращает названия имён файла и удаляет их из базы данных
func (i *ImageDB) DeleteImage(id string) (bool, string, string, error) {
	var file, preview string
	row := i.db.QueryRow("select name, preview_name from images where id=$1", id)
	err := row.Scan(&file, &preview)
	if err != nil {
		return false, file, preview, err
	}
	_, err = i.db.Exec("DELETE FROM images WHERE id=$1", id)
	if err != nil {
		return false, file, preview, err
	}
	return true, file, preview, err
}

// AddImage добавляет запись с именем файла и именем файла-миниатюры в базу данных
func (i *ImageDB) AddImage(fileName string) (bool, error) {
	id := utils.RandStringRunes(5)
	_, err := i.db.Exec("insert into images (id, name, preview_name) values ($1,$2,$3)", id, fileName, "pre"+fileName)
	if err != nil {
		return false, err
	}
	return true, err
}
