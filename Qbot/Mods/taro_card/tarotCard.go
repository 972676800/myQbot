package taro_card

import (
	"Qbot_gocode/db"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"io/ioutil"
	"math/rand"
)

var TaroCard Tc

type Tc map[int]Image

type Image struct {
	gorm.Model
	Id         uint
	URL        string `yaml:"url"`
	Name       string `yaml:"name"`
	ProductID  uint
	Product    Product           `yaml:"position"`
	ProductMap map[string]string `gorm:"json"`
}

func (Image) TableName() string {
	return "t_taro_card"
}

type Product struct {
	gorm.Model
	Upright  string `yaml:"upright"`
	Reversed string `yaml:"reversed"`
}

func (Product) TableName() string {
	return "t_taro_product"
}

func (t *Tc) Select() (string, string) {
	for _, v := range *t {
		var minging string
		var position string
		for k, v1 := range v.ProductMap {
			if rand.Intn(10) < 7 {
				position = v1
				minging = k
				break
			}
			position = v1
			minging = k
			break
		}
		return v.URL, fmt.Sprintf("您抽到的是【%s】【%s】\n%s", position, v.Name, minging)
	}
	return "", "发生错误，请联系主人"
}

func GetTaroCardFromYml() {
	yamlFile, err := ioutil.ReadFile("tlp.yml")
	if err != nil {
		panic(err)
	}

	var images []Image
	err = yaml.Unmarshal(yamlFile, &images)
	if err != nil {
		panic(err)
	}

	TaroCard = make(Tc, len(images))
	for k, image := range images {
		image.ProductMap = make(map[string]string, 2)
		TaroCard[k] = image
		TaroCard[k].ProductMap[image.Product.Upright] = "正位"
		TaroCard[k].ProductMap[image.Product.Reversed] = "逆位"
		fmt.Printf("Name: %s\nURL: %s\n", image.Name, image.URL)
		//for _, pos := range image.Product {
		//	fmt.Printf("Upright: %s\n", pos.Upright)
		//	fmt.Printf("Reversed: %s\n", pos.Reversed)
		//}
		fmt.Println()
	}

	e := db.DB.AutoMigrate(&Image{})
	if e != nil {
		println(e.Error())
	}
	SaveTaroCard()
}

func SaveTaroCard() error {
	for _, image := range TaroCard {
		if err := db.DB.Create(&image).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetImagesFromDatabase() error {
	var images []Image
	if err := db.DB.Find(&images).Error; err != nil {
		return err
	}
	TaroCard = make(Tc, len(images))
	for _, image := range images {
		TaroCard[int(image.Id)] = image
		//TaroCard[int(image.Id)].ProductMap[image.Product.Upright] = "正位"
		//TaroCard[int(image.Id)].ProductMap[image.Product.Reversed] = "逆位"
	}
	fmt.Printf("%v", TaroCard)

	return nil
}

func init() {
	GetImagesFromDatabase()
	fmt.Printf("%v", TaroCard)
}
