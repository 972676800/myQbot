package user

import (
	"Qbot_gocode/Mods/taro_card"
	"Qbot_gocode/db"
	"encoding/json"
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var PalyerList []Player

type Player struct {
	gorm.Model
	UserId   int `gorm:"uniqueIndex"`
	NickName string
	Bag      *Bag  `gorm:"json"`
	Point    *Sign `gorm:"json"`
}

func (Player) TableName() string {
	return "t_player"
}

func newPlayer(userId int, nickName string) *Player {
	var tmp *Bag
	tmp = new(Bag)
	return &Player{
		UserId:   userId,
		NickName: nickName,
		Bag:      tmp,
		Point:    new(Sign),
	}
}

func FindPlayer(ctx *zero.Ctx) *Player {
	userId := int(ctx.Event.Sender.ID)
	userName := fmt.Sprintf("伟大的-%v-%v-%v", ctx.Event.Sender.Title, ctx.Event.Sender.Card, ctx.Event.Sender.NickName)
	for _, v := range PalyerList {
		if v.UserId == userId {
			v.NickName = userName
			return &v
		}
	}
	p := newPlayer(userId, userName)
	PalyerList = append(PalyerList, *p)
	return p
}

func (p *Player) DefaultHandler(ctx *zero.Ctx) {
	//全文匹配获取指令
	for _, msg := range ctx.Event.Message {
		switch msg.Type {
		case "text":
			if _, ok := Commands[msg.Data["text"]]; ok {
				cmd := msg.Data["text"]
				if cmd == "签到" {
					p.Sign(ctx)
					return
				}
				if cmd == "查看商店" {
					Things.show()
					ctx.SendChain(message.Text(fmt.Sprintf("商店\n(物品名称 所需积分)：\n%v", Things.show())))
				}
				if cmd == "查看背包" {
					p.Bag.showBags()
					ctx.SendChain(message.Text(fmt.Sprintf("【%s】 当前背包：\n%v", p.NickName, p.Bag.showBags())))
				}
				if cmd == "查看积分" {
					ctx.SendChain(message.Text(fmt.Sprintf("【%s】\n当前积分：%v", p.NickName, p.Point.showPoints())))
				}
				if cmd == "抽塔罗牌" {
					url, str := taro_card.TaroCard.Select()
					ctx.SendChain(message.Image(url))
					time.Sleep(time.Second * 1)
					ctx.SendChain(message.Text(fmt.Sprintf("【%s】\n%s", p.NickName, str)))
				}
				if cmd == "使用指南" {
					var result string
					for k, v := range commands_slice {
						if v != "使用指南" {
							result += fmt.Sprintf("%d、%s\n", k+1, v)
						}
					}
					ctx.SendChain(message.Text(fmt.Sprintf("可用指令：\n%v", result)))
				}
			}
			if strings.Contains(msg.Data["text"], "捏捏牛牛") {
				p.Bag.NNN(ctx, contain(msg.Data["text"]), p)
			}
			if strings.Contains(msg.Data["text"], "购买") {
				p.Buy(ctx, msg)
			}
		default:
		}
	}
}

func contain(s string) int {
	for i := 0; i < 10; i++ {
		if strings.Contains(s, strconv.Itoa(i+1)) {
			return i
		}
	}
	return -1
}
func (p *Player) Sign(ctx *zero.Ctx) {

	if p.Point.isSigned() {
		ctx.SendChain(message.Text(fmt.Sprintf("【%s】 签到成功 积分+%d\n当前积分:%d", p.NickName, p.Point.pointAdd(), p.Point.Points)))
		return
	}
	ctx.SendChain(message.Text(fmt.Sprintf("【%s】 今日已签到\n当前积分:%d", p.NickName, p.Point.Points)))
}

func (p *Player) Buy(ctx *zero.Ctx, msg message.MessageSegment) {
	if price, ok := Things[strings.Split(msg.Data["text"], " ")[1]]; ok {
		name := strings.Split(msg.Data["text"], " ")[1]
		if !Things.IsApplied(name, p.Point.Points) {
			ctx.SendChain(message.Text(fmt.Sprintf("【%s】 购买【%v】失败\n所需积分%d\n当前积分:%d\n", p.NickName, name, price, p.Point.Points)))
			return
		}
		p.Point.buy(price)
		p.Bag.buy(name)
		ctx.SendChain(message.Text(fmt.Sprintf("【%s】 购买【%v】成功\n消费积分：%d\n当前积分:%d\n当前背包：\n%v", p.NickName, name, price, p.Point.Points, p.Bag.showBags())))
		return
		//p.Buy(ctx, strings.Split(msg.Data["text"], " ")[1], p.Point.Points)
	}
	ctx.SendChain(message.Text("没有此类物品"))

}

func init() {
	e := db.DB.AutoMigrate(&Player{})
	if e != nil {
		println(e.Error())
	}

}

func GetPlayersFromDatabase() {
	// 2. 查询数据库中的Player记录
	var players []Player
	if err := db.DB.Find(&players).Error; err != nil {
		panic(err)
	}
	// 3. 转换查询结果为Player对象并添加到palyerList中
	for _, p := range players {
		b, err := json.Marshal(p.Bag)
		if err != nil {
			panic(err)
		}
		var tmp Bag
		err = json.Unmarshal(b, &tmp)
		if err != nil {
			panic(err)
		}
		p.Bag = &tmp
		s, err := json.Marshal(p.Point)
		if err != nil {
			panic(err)
		}
		var signTmp Sign
		err = json.Unmarshal(s, &signTmp)
		if err != nil {
			panic(err)
		}
		p.Point = &signTmp
		PalyerList = append(PalyerList, p)
	}
}

func SavePlayers() {

	for _, player := range PalyerList {
		var existingPlayer Player
		if err := db.DB.Where("user_id = ?", player.UserId).First(&existingPlayer).Error; err != nil {
			// 如果查不到记录，则创建新的记录
			if err := db.DB.Create(player).Error; err != nil {
				fmt.Printf("failed to create player record: %v\n", err)
			}
		} else {
			// 如果查到了记录，则更新该记录
			if err := db.DB.Save(player).Error; err != nil {
				fmt.Printf("failed to update player record: %v\n", err)
			}
		}
	}
}

func init() {
	GetPlayersFromDatabase()
}
