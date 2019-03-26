package modules

import (
	"fmt"
        "strings"
        "net/http"
	"encoding/json"
)

const (
	weatherCacheKey = "weather:location"
	OpenWeatherMapToken = "af7bfe22287c652d032a3064ffa44088"
	snowIcon            = "\uf2dc"
	rainIcon            =  "\uf73d" // "\uf740"
	sunIcon             = "\uf185"
	cloudIcon           = "\uf0c2"
	thunderIcon         = ""

	smogIcon    = "\uf75f"
	celsiusIcon = "\u00B0" //"\u2103"
)

var (
	iconSet = map[string]string{"Clear": sunIcon, "Clouds": cloudIcon, "Thunderstorm": thunderIcon,
		"Drizzle": rainIcon, "Rain": rainIcon, "Snow": snowIcon}
)

type location struct {
        Name    string `json:"city"`
        Coord   string `json:"loc"`
        rewrite string
}

type weather struct {
        Weather []struct {
                Main string `json:"main"`
                Desc string `json:"description"`
        } `json:"weather"`
        Main struct {
                Temp float64 `json:"temp"`
        } `json:"main"`
        Wind struct {
                Speed int `json:"speed"`
        } `json:"wind"`
}

type Weather struct {
        name string
}


func getWeatherModule(mo *ModuleOutput, cfg ModuleConfig) {
	var loc location
	v := cache.Get(weatherCacheKey)
	if v != nil {
		loc, _ = v.(location)
	} else {
		//module just started
		_, ok := cfg.Extra["location"] 
		if ok {
			city, _ := cfg.Extra["location"].(string)
			if city == "auto" {
				loc = getLocation()
			} else {
				loc.rewrite = city
				loc.Name = city
			}
		} else {
			loc = getLocation()
		}
		cache.Add(weatherCacheKey, loc, "24h")
	}
	wf := getWeather(loc)
	icon, ok := iconSet[wf.Weather[0].Main]
	if !ok {
		icon = smogIcon
	}
	mo.FullText = fmt.Sprintf("%s%s: %s %.0f%s (%d m/s)", mo.FullText, loc.Name, icon, wf.Main.Temp, celsiusIcon, wf.Wind.Speed)
	
}

func getLocation() location {
	loc := new(location)
	url := "http://ipinfo.io/json"
	res, err := http.Get(url)
	if err != nil {
		return *loc
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(loc); err != nil {
		panic(err)
	}
	return *loc
}

func getWeather(l location) *weather {
	w := new(weather)
	coord := strings.Split(l.Coord, ",")
	var url string
	if len(l.rewrite) > 0 {
		url = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s&units=metric", OpenWeatherMapToken, l.rewrite)
	} else {
		url = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&lat=%s&lon=%s&units=metric", OpenWeatherMapToken, coord[0], coord[1])
	}
	res, err := http.Get(url)
	if err != nil {
		return w
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(w); err != nil {
		return w
	}
	return w

}

func init() {
        RegisteredFuncs["weather"] = getWeatherModule
}
