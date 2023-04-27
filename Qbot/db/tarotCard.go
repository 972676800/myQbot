package db

import (
	"Qbot_gocode/Mods/taro_card"
	"fmt"
)

func GetImagesFromDatabase() error {
	var images []taro_card.Image
	if err := DB.Find(&images).Error; err != nil {
		return err
	}
	taro_card.TaroCard = make(taro_card.Tc, len(images))
	for _, image := range images {
		taro_card.TaroCard[int(image.Id)] = image
		//TaroCard[int(image.Id)].ProductMap[image.Product.Upright] = "正位"
		//TaroCard[int(image.Id)].ProductMap[image.Product.Reversed] = "逆位"
	}
	fmt.Printf("%v", taro_card.TaroCard)

	return nil
}
