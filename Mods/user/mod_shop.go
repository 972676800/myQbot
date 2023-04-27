package user

import (
	"Qbot_gocode/db"
	"fmt"
	"gorm.io/gorm"
)

var Things things

type things map[string]int

type thing struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex"`
	Price int
}

func (thing) TableName() string {
	return "t_shop"
}

func (t things) IsApplied(s string, i int) bool {
	if t[s] <= i {
		return true
	}
	return false
}

func (t things) show() string {
	var result string
	for k, v := range t {
		result += fmt.Sprintf("【%s】:%d\n", k, v)
	}
	return result
}

func (t things) buy(s string, i int) bool {
	if t[s] <= i {
		return true
	}
	return false
}

func init() {

	e := db.DB.AutoMigrate(&thing{})
	if e != nil {
		println(e.Error())
	}
	GetThingsFromDataBase()
	fmt.Printf("%v", Things)
}

func GetThingsFromDataBase() {
	// 查询数据
	var thingList []thing
	if err := db.DB.Find(&thingList).Error; err != nil {
		fmt.Printf("Failed to get things records: %v\n", err)
	}

	// 处理数据
	Things = make(things)
	for _, item := range thingList {
		Things[item.Name] = item.Price
	}
}
