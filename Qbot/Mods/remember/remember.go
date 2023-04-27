package remember

import (
	"Qbot_gocode/Mods/utils"
	"Qbot_gocode/db"
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
	"math/rand"
	"regexp"
)

var R remember

type remember struct {
	gorm.Model
	Id uint
	M  map[string]*rememb `gorm:"json"`
}
type rememb struct {
	Str string
	Ops []string
}

func (remember) TableName() string {
	return "t_remember"
}

func (r *remember) SaveToDB() error {
	err := db.DB.Exec("DELETE FROM t_remember").Error
	if err != nil {
		return err
	}

	err = db.DB.Create(&R).Error
	if err != nil {
		return (err)
	}
	return nil
}

func FindRemember() *remember {
	if R.M == nil {
		R.M = make(map[string]*rememb)
	}

	return &R
}

func (r *remember) DefaultHandler(ctx *zero.Ctx) {
	if s, ok := utils.FindKV(ctx.MessageString()); ok {
		switch s[0] {
		case "记住":
			if s[2] != "" {
				r.M[s[1]] = &rememb{Str: s[2], Ops: nil}
				ctx.SendChain(message.Text(fmt.Sprintf("【%v】记住了！", s[1])))
			} else {
				ctx.SendChain(message.Text(fmt.Sprintf("【%v】内容丢失，【记住】+【关键字】+【内容】中间要有空格哦亲~！", s[1])))
			}

			return
		case "添加":
			if s[2] != "" {
				if _, ok := r.M[s[1]]; !ok {
					r.M[s[1]] = &rememb{Str: ""}
				}
				r.M[s[1]].Ops = append(r.M[s[1]].Ops, s[2])
				ctx.SendChain(message.Text(fmt.Sprintf("【%v】添加了！", s[1])))
			} else {
				ctx.SendChain(message.Text(fmt.Sprintf("【%v】内容丢失，【记住】+【关键字】+【内容】中间要有空格哦亲~！", s[1])))
			}
			return

		case "查看":
			if s[1] != "" {
				if str, ok := r.M[s[1]]; ok {
					for _, v := range r.FindAllValue(str) {
						if url, yes := IsStrUrl(v); yes {
							ctx.SendChain(message.Image(url))
						} else {
							ctx.SendChain(message.Text(url))
						}
					}

				}
			}
		case "忘掉":
			if s[1] != "" {
				delete(r.M, s[1])
				ctx.SendChain(message.Text(fmt.Sprintf("【%v】忘掉了！", s[1])))
			}
		}
	}

	if str, ok := r.M[ctx.MessageString()]; ok {
		if url, yes := IsStrUrl(r.FindOneValue(str)); yes {
			ctx.SendChain(message.Image(url))
		} else {
			ctx.SendChain(message.Text(url))
		}
	}
}

func IsStrUrl(str string) (string, bool) {
	re := regexp.MustCompile(`url=([^&]+)`)
	result := re.FindStringSubmatch(str)
	if len(result) < 2 {
		return str, false
	} else {
		return result[1], true
	}
}

func (r *remember) FindOneValue(str *rememb) string {
	if str.Ops != nil {
		return str.Ops[rand.Intn(len(str.Ops))]
	}
	return str.Str
}

func (r *remember) FindAllValue(str *rememb) []string {
	var result []string
	if str.Ops != nil {
		for _, op := range str.Ops {
			result = append(result, op)
		}

	}
	result = append(result, str.Str)
	return result
}

func init() {
	e := db.DB.AutoMigrate(&remember{})
	if e != nil {
		println(e.Error())
	}

}
