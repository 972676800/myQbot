package main

import (
	"Qbot_gocode/Mods/fd"
	"Qbot_gocode/Mods/reverse"
	"Qbot_gocode/Mods/user"
	"Qbot_gocode/mods/remember"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := zero.New()
	engine.OnFullMatch("bot reverse").Handle(
		func(ctx *zero.Ctx) {
			remember.R.SaveToDB()
			user.SavePlayers()
			ctx.SendChain(message.Text("reverse success"))
			//ctx.SendChain(message.Text("hello world!"))
		})

	// 天气API
	engine.OnMessage().Handle(
		func(ctx *zero.Ctx) {
			//go weather.DefaultHandle(ctx)
			//ctx.SendChain(message.Text("hello world!"))
		})
	engine.OnPrefix("hello").Handle(
		func(ctx *zero.Ctx) {
			//ctx.SendChain(message.Text("hello world!"))
		})

	// remember复用
	engine.OnMessage().Handle(
		func(ctx *zero.Ctx) {
			go remember.FindRemember().DefaultHandler(ctx)
		})
	// 主要
	engine.OnMessage().Handle(
		func(ctx *zero.Ctx) {
			go user.FindPlayer(ctx).DefaultHandler(ctx)
			//ctx.SendChain(message.Text("hello world!"))
		})

	// 复读
	engine.OnMessage().Handle(
		func(ctx *zero.Ctx) {
			go reverse.NewReverse().ReverseHandler(ctx)
		})

	// 吃什么
	engine.OnFullMatch("吃什么").Handle(
		func(ctx *zero.Ctx) {
			go fd.RandomFood(ctx)
		})

	engine.OnFullMatch("小黑子").Handle(
		func(ctx *zero.Ctx) {
			ctx.SendChain(message.Image("https://s1.ax1x.com/2023/04/11/ppLieit.png"))
		})
}
