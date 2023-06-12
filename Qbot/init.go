package main

import (
	"Qbot_gocode/Mods/API"
	"Qbot_gocode/Mods/fd"
	"Qbot_gocode/Mods/remember"
	"Qbot_gocode/Mods/reverse"
	"Qbot_gocode/Mods/user"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

func init() {
	engine := zero.New()
	go func() {
		for {
			user.SavePlayers()
			remember.Save()
			time.Sleep(time.Hour * 1)
		}
	}()

	// 天气API
	engine.OnMessage().Handle(
		func(ctx *zero.Ctx) {
			//go weather.DefaultHandle(ctx)
			//ctx.SendChain(message.Text("hello world!"))
		})

	// ChatGPT
	engine.OnPrefix("g").Handle(
		func(ctx *zero.Ctx) {
			a := ctx.MessageString()[1:]

			gptClient, ok := API.GptClient[int(ctx.Event.GroupID)]
			if !ok {
				gptClient = API.NewGptClient()
				API.GptClient[int(ctx.Event.GroupID)] = gptClient
			}
			ctx.SendChain(message.Text(gptClient.DefaultHandler(a)))
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
	engine.OnKeyword("想吃").Handle(
		func(ctx *zero.Ctx) {
			go fd.InsertFood(ctx)
		})

	engine.OnFullMatch("小黑子").Handle(
		func(ctx *zero.Ctx) {
			ctx.SendChain(message.Image("https://s1.ax1x.com/2023/04/11/ppLieit.png"))
		})
	engine.OnKeyword("哭哭").Handle(
		func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("哭你妈 再哭给你两拳"))
		})
}
