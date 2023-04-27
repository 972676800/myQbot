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
	M  map[string]*Rememb `gorm:"json"`
}

type Rememb struct {
	Str string
	Ops []string
}

type r struct {
	Id   uint
	Name string
	Str  string
	Ops  []string `gorm:"json"`
}

func (r) TableName() string {
	return "t_remembers"
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

func Save() {
	var temp []r
	for K, V := range R.M {
		t := r{
			Name: K,
			Str:  V.Str,
			Ops:  V.Ops,
		}
		temp = append(temp, t)
	}
	for _, v := range temp {
		var t r
		if err := db.DB.Where("name = ?", v.Name).First(&t).Error; err != nil {
			if err := db.DB.Create(&v).Error; err != nil {
				fmt.Printf("failed to create player record: %v\n", err)
			}
		} else {
			// 如果查到了记录，则更新该记录
			if err := db.DB.Save(&v).Error; err != nil {
				fmt.Printf("failed to update player record: %v\n", err)
			}
		}
	}
}

func FindRemember() *remember {
	if R.M == nil {
		R.M = make(map[string]*Rememb)
	}

	return &R
}

func (r *remember) DefaultHandler(ctx *zero.Ctx) {
	if s, ok := utils.FindKV(ctx.MessageString()); ok {
		switch s[0] {
		case "记住":
			if s[2] != "" {
				r.M[s[1]] = &Rememb{Str: s[2], Ops: nil}
				ctx.SendChain(message.Text(fmt.Sprintf("【%v】记住了！", s[1])))
			} else {
				ctx.SendChain(message.Text(fmt.Sprintf("【%v】内容丢失，【记住】+【关键字】+【内容】中间要有空格哦亲~！", s[1])))
			}

			return
		case "添加":
			if s[2] != "" {
				if _, ok := r.M[s[1]]; !ok {
					r.M[s[1]] = &Rememb{Str: ""}
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

func (r *remember) FindOneValue(str *Rememb) string {
	if str.Ops != nil {
		return str.Ops[rand.Intn(len(str.Ops))]
	}
	return str.Str
}

func (r *remember) FindAllValue(str *Rememb) []string {
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

	e1 := db.DB.AutoMigrate(&r{})
	if e1 != nil {
		println(e1.Error())
	}
	//GetRememberFromDatabase()
	f()
}

func f() {
	var t []r
	if err := db.DB.Find(&t).Error; err != nil {
		panic(err)
	}
	R.M = make(map[string]*Rememb, len(t))
	for _, v := range t {
		R.Id = v.Id
		R.M[v.Name] = &Rememb{
			Str: v.Str,
			Ops: v.Ops,
		}
	}
	fmt.Printf("%v", R)
}

func GetRememberFromDatabase() error {
	err := db.DB.Table(R.TableName()).First(&R).Error
	if err != nil {
		return err
	}
	return nil
}
