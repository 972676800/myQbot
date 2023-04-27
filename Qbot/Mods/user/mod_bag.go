package user

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

type bag struct {
	Things []map[string]*value `gorm:"json"`
}

type value struct {
	V      int
	TimeOK bool
	LastDo time.Time
}

func (b *bag) showBags() string {
	var result string
	for index, value := range b.Things {
		for k, v := range value {
			result += fmt.Sprintf("%d:【%s】成长值【%d】\n", index+1, k, v.V)
		}
	}
	return result
}

func (b *bag) buy(n string) {
	tmp := make(map[string]*value, 1)
	tmp[n] = &value{
		V:      1,
		TimeOK: false,
	}
	//tmp[n]++
	b.Things = append(b.Things, tmp)
}

func (b *bag) NNN(ctx *zero.Ctx, index int, p *Player) {
	msgIndex := index + 1
	for i, v := range b.Things {
		if index == i {
			if v["牛牛"].isTimeOk() {
				v["牛牛"].LastDo = time.Now()
				v["牛牛"].V++
				ctx.SendChain(message.Text(fmt.Sprintf("【%s】\n当前【%d】号牛牛【成长值】：%V", p.NickName, msgIndex, b.Things[index]["牛牛"].V)))
				return
			} else {
				ctx.SendChain(message.Text(fmt.Sprintf("【%s】\n你今天已经捏过【%d】号牛牛啦~\n当前【成长值】：%V", p.NickName, msgIndex, b.Things[index]["牛牛"].V)))
				return
			}
		}
	}
	ctx.SendChain(message.Text(fmt.Sprintf("【%s】\n当前没有【%d】号牛牛", p.NickName, msgIndex)))
}
func (v *value) isTimeOk() bool {
	//获得上次操作时间 转化为整点
	LastDoTimeStamp := v.LastDo.Truncate(24 * time.Hour)
	//获取现在时间 转化为整点
	now := time.Now()
	zeroClock := now.Truncate(24 * time.Hour)
	if LastDoTimeStamp == zeroClock {
		return false
	}
	return true
}
