package reverse

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var Reverse *reverse

type reverse struct {
	isReverse   bool
	msgRemember []string
}

func (r *reverse) ReverseHandler(ctx *zero.Ctx) {
	for _, msg := range ctx.Event.Message {
		switch msg.Type {
		case "text":
			LeftPush(msg.Data["text"])
			switch msg.Data["text"] {
			case "开启复读":
				r.openReverse()
				ctx.SendChain(message.Text("开启成功！"))
			case "关闭复读":
				r.closeReverse()
				ctx.SendChain(message.Text("关闭成功！"))
			}
		}
	}
	if r.isReverse {
		if r.msgRemember[0] == r.msgRemember[1] {
			ctx.SendChain(message.Text(r.msgRemember[0]))
		}
	}
}

func LeftPush(s string) {
	Reverse.msgRemember = Reverse.msgRemember[1:]
	Reverse.msgRemember = append(Reverse.msgRemember, s)
}

func NewReverse() *reverse {
	if Reverse == nil {
		Reverse = &reverse{isReverse: false, msgRemember: make([]string, 2, 2)}
	}
	return Reverse
}

func (r *reverse) openReverse() {
	r.isReverse = true
}

func (r *reverse) closeReverse() {
	r.isReverse = false
}

func init() {
	Reverse = NewReverse()
}
