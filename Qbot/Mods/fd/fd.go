package fd

import (
	"Qbot_gocode/db"
	"errors"
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"regexp"
)

var F []Food

type Food struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

func (Food) TableName() string {
	return "t_food"
}

func init() {
	// 使用AutoMigrate方法自动创建数据表
	db.DB.AutoMigrate(&Food{})
	db.DB.Find(&F)

}

func RandomFood(ctx *zero.Ctx) {
	ctx.SendChain(message.Text("吃" + F[rand.Intn(len(F))].Name))
}

func InsertFood(ctx *zero.Ctx) {
	name, err := FindString(ctx.MessageString())
	if err != nil && name != "" {
		ctx.SendChain(message.Text(err.Error()))
		return
	}
	var existingFood Food
	if err := db.DB.Where("name = ?", name).First(&existingFood).Error; err != nil {
		newFood := Food{
			ID:   getMax() + 1,
			Name: name,
		}
		// 如果查不到记录，则创建新的记录
		if err := db.DB.Create(&newFood).Error; err != nil {
			fmt.Printf("failed to create player record: %v\n", err)
		}
		F = append(F, newFood)
	} else {
		// 如果查到了记录，则更新该记录
		ctx.SendChain(message.Text("菜单已存在：" + name))
	}

	db.DB.Find(&F)
	ctx.SendChain(message.Text("菜单添加了：" + name))
}

func getMax() int {
	temp := 0
	for _, food := range F {
		if food.ID > temp {
			temp = food.ID
		}
	}
	return temp
}

func FindString(s string) (string, error) {
	// 定义正则表达式
	reg := regexp.MustCompile(`想吃 (.*)`)
	// 匹配字符串
	res := reg.FindStringSubmatch(s)
	// 输出匹配结果
	if res != nil {
		return res[1], nil
	}
	return "", errors.New("添加失败，格式为【想吃】+【空格】+【名字】")
}
