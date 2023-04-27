package db

import (
	"Qbot_gocode/Mods/user"
	"encoding/json"
)

func GetPlayersFromDatabase() {
	// 2. 查询数据库中的Player记录
	var players []user.Player
	if err := DB.Find(&players).Error; err != nil {
		panic(err)
	}
	// 3. 转换查询结果为Player对象并添加到palyerList中
	for _, p := range players {
		b, err := json.Marshal(p.Bag)
		if err != nil {
			panic(err)
		}
		var tmp user.Bag
		err = json.Unmarshal(b, &tmp)
		if err != nil {
			panic(err)
		}
		p.Bag = &tmp
		s, err := json.Marshal(p.Point)
		if err != nil {
			panic(err)
		}
		var signTmp user.Sign
		err = json.Unmarshal(s, &signTmp)
		if err != nil {
			panic(err)
		}
		p.Point = &signTmp
		user.PalyerList = append(user.PalyerList, p)
	}
}
