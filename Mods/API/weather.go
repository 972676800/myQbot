package weather

import (
	"encoding/json"
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const (
	MSN = "https://restapi.amap.com/v3/weather/weatherInfo"
	KEY = "381a7c1d827d75bfd5a141468d46821d"
)

type Weather struct {
	Data     *Data
	Url      string
	Day      string
	City     string
	CityCode int
}

type Data struct {
	Status    string `json:"status"`
	Count     string `json:"count"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Forecasts []struct {
		City       string `json:"city"`
		Adcode     string `json:"adcode"`
		Province   string `json:"province"`
		Reporttime string `json:"reporttime"`
		Casts      []struct {
			Date           string `json:"date"`
			Week           string `json:"week"`
			Dayweather     string `json:"dayweather"`
			Nightweather   string `json:"nightweather"`
			Daytemp        string `json:"daytemp"`
			Nighttemp      string `json:"nighttemp"`
			Daywind        string `json:"daywind"`
			Nightwind      string `json:"nightwind"`
			Daypower       string `json:"daypower"`
			Nightpower     string `json:"nightpower"`
			DaytempFloat   string `json:"daytemp_float"`
			NighttempFloat string `json:"nighttemp_float"`
		} `json:"casts"`
	} `json:"forecasts"`
}

func DefaultHandle(ctx *zero.Ctx) (w *Weather) {
	raw := strings.TrimSpace(ctx.MessageString())
	re := regexp.MustCompile(`查询(\p{Han}{1,9}[市区])?(.+)?天气`)

	if re.MatchString(raw) {
		w = new(Weather)
		value := re.FindStringSubmatch(raw)
		w.Init(value)
		ctx.SendChain(message.Text(w.Result()))
	}
	return
}

func (w *Weather) Init(raw []string) {
	w.getCityCode(raw[1])
	w.getDay(raw[2])
	w.setUrl()
	w.getData()

}

func (w *Weather) setUrl() {
	w.Url = fmt.Sprintf("%s?key=%s&city=%d&extensions=all", MSN, KEY, w.CityCode)
}

func (w *Weather) getData() {
	w.Data = new(Data)
	r, _ := http.NewRequest("GET", w.Url, nil)
	re, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	raw, _ := io.ReadAll(re.Body)
	json.Unmarshal(raw, &w.Data)
	return
}

func (w *Weather) getDay(raw string) {
	if strings.Contains("今日明日近期", raw) {
		w.Day = raw
	} else {
		w.Day = "今日"
	}
}

func (w *Weather) getCityCode(raw string) {
	w.CityCode = 110000
}

func (w *Weather) isOk() bool {
	return w.Data.Status == "1"
}

func (w Weather) Result() string {
	if w.isOk() {
		switch w.Day {
		case "今日":
			return fmt.Sprintf("%v", w.Data.Forecasts[0].Casts[0])
		case "明日":
			return fmt.Sprintf("%v", w.Data.Forecasts[1].Casts[0])
		case "近期":
			return fmt.Sprintf("%v", w.Data.Forecasts[1].Casts)
		}
	} else {
		return fmt.Sprintf("不知道也做不到")
	}
	return ""
}
