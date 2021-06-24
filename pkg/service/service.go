package service

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"gora/pkg/db"
	"gora/pkg/utils"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/nfnt/resize"
)

const (
	baseImagePath = "img/"
)

type ImageService struct {
	db *db.SQLiteDB
}

func NewImageService(db *db.SQLiteDB) *ImageService {
	return &ImageService{db}
}

type ImageServiceInterfase interface {
	GetImage(id string) (string, error)
	AddImage(io.Reader) bool
	GetImageList() ([]byte, error)
	DeleteImage(id string) bool
}

type Service struct {
	ImageServiceInterfase
}

func NewService(db *db.SQLiteDB) *Service {
	return &Service{
		ImageServiceInterfase: NewImageService(db),
	}
}

//GetImagePreview возвращаеь файл преобразованный в строку base64
func (i *ImageService) GetImage(id string) (string, error) {

	filename, err := i.db.GetImagePreview(id)
	var encdata string
	if err != nil {
		log.Println("filename generate error", err)
		return encdata, err
	}
	file, err := os.Open(baseImagePath + filename)
	if err != nil {
		return encdata, err
	}
	defer file.Close()
	r := bufio.NewReader(file)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("filename ", err)
		return encdata, err
	}
	encdata = base64.StdEncoding.EncodeToString(data)

	return encdata, nil
}

// AddImage сохряняет файл в папку, генерирует и сохряняет миниатюру
func (i *ImageService) AddImage(data io.Reader) bool {

	fileName := utils.RandStringRunes(10) + ".jpeg"
	file, err := os.Create(baseImagePath + fileName)
	if err != nil {
		log.Println("file error", err)
		return false
	}
	io.Copy(file, data)
	file.Close()

	previewfile, err := os.Create(baseImagePath + "pre" + fileName)
	if err != nil {
		os.Remove(baseImagePath + fileName)
		log.Println("file error", err)
		return false
	}
	defer previewfile.Close()
	file, err = os.Open(baseImagePath + fileName)
	if err != nil {
		os.Remove(baseImagePath + fileName)
		os.Remove(baseImagePath + "pre" + fileName)
		log.Println("file error", err)
		return false
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		os.Remove(baseImagePath + fileName)
		os.Remove(baseImagePath + "pre" + fileName)
		log.Println(err)
		return false
	}
	defer file.Close()
	err = jpeg.Encode(previewfile, resize.Resize(100, 70, img, resize.NearestNeighbor), nil)
	if err != nil {
		os.Remove(baseImagePath + fileName)
		os.Remove(baseImagePath + "pre" + fileName)
		log.Println(err)
		return false
	}
	result, err := i.db.AddImage(fileName)
	if err != nil {
		os.Remove(baseImagePath + fileName)
		os.Remove(baseImagePath + "pre" + fileName)
		log.Println(err)
	}

	return result
}

// DeleteImage удаляет сохранённые файлы
func (i *ImageService) DeleteImage(id string) bool {
	resp, file, preview, err := i.db.DeleteImage(id)
	if err != nil {
		log.Panicln(err)
	}
	os.Remove(baseImagePath + file)
	os.Remove(baseImagePath + preview)
	return resp
}

func (i *ImageService) GetImageList() ([]byte, error) {
	imgs, err := i.db.GetImages()
	if err != nil {
		return nil, err
	}

	for i, img := range imgs.Images {
		file, err := os.Open(baseImagePath + img.PeviewFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		r := bufio.NewReader(file)
		data, err := ioutil.ReadAll(r)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		img.Preview = base64.StdEncoding.EncodeToString(data)
		imgs.Images[i] = img
	}
	return json.Marshal(imgs)

}
