package fd

import (
	"Qbot_gocode/db"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
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
