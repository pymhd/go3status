package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	weatherCacheKey     = "weather:location"
	OpenWeatherMapToken = "af7bfe22287c652d032a3064ffa44088"
	snowIcon            = "\uf2dc"
	rainIcon            = "\uf73d" // "\uf740"
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
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

func getWeatherModule(mo *ModuleOutput, cfg ModuleConfig) {
	log.Debug(`"Weather" module hook started`)
	var loc location
	v := cache.Get(weatherCacheKey)
	if v != nil {
		loc, _ = v.(location)
		log.Debugf("Found location in cache: %s\n", loc.Name)
	} else {
		//module just started
		log.Debug("Could not find location in cache, config lookup...")
		_, ok := cfg.Extra["location"]
		if ok {
			city, _ := cfg.Extra["location"].(string)
			if city == "auto" {
				log.Debug("Found auto detect location request in config")
				loc = getLocation()
			} else {
				loc.rewrite = city
				loc.Name = city
				log.Debugf("Found location in config: %s\n", loc.Name)
			}
		} else {
			log.Debug("Could not find location in config, will auto detect it by ip addr using ipinfo.io svc")
			loc = getLocation()
		}
		cache.Add(weatherCacheKey, loc, "24h")
		log.Debug("Location pushed in cache")
	}
	wf := getWeather(loc)
	if wf == nil {
		log.Debug("Could not fetch weather forecast for openweathermap API. Returning N/A value")
		mo.FullText += "N/A"
		return
	}
	icon, ok := iconSet[wf.Weather[0].Main]
	if !ok {
		icon = smogIcon
	}
	if cfg.ShortFormat {
		mo.FullText = fmt.Sprintf("%s", icon)
	} else {
		mo.FullText = fmt.Sprintf("%s%s: %s %.0f%s (%.1f m/s)", mo.FullText, loc.Name, icon, wf.Main.Temp, celsiusIcon, wf.Wind.Speed)
	}
	log.Debugf("Returning -> %s%s: %s %.0f%s (%.1f m/s) \n", mo.FullText, loc.Name, icon, wf.Main.Temp, celsiusIcon, wf.Wind.Speed)

}

func getLocation() location {
	log.Debug("Location auto detect started")
	loc := new(location)
	url := "http://ipinfo.io/json"
	res, err := http.Get(url)
	if err != nil {
		log.Errorf("Could not reach ipinfo.io svc: (%s)\n", err)
		return *loc
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(loc); err != nil {
		log.Errorf("Could not parse json obj in response from ipinfo.io (%s)\n", err)
		return *loc
	}
	log.Debug("Successfully found location by ip addr")
	return *loc
}

func getWeather(l location) *weather {
	log.Debug("Weather forecast fetching started")
	w := new(weather)
	coord := strings.Split(l.Coord, ",")
	var url string
	if len(l.rewrite) > 0 {
		log.Debugf("Going to get weather for manually specified region: %s\n", l.rewrite)
		url = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s&units=metric", OpenWeatherMapToken, l.rewrite)
	} else {
		log.Debug("Going to get weather for auto detected region by its coords")
		url = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&lat=%s&lon=%s&units=metric", OpenWeatherMapToken, coord[0], coord[1])
	}
	res, err := http.Get(url)
	if err != nil {
		log.Errorf("Could not fetch weather forcast: (%s)\n", err)
		return nil
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(w); err != nil {
		log.Errorf("Could not parse json obj in response from openweathermap API (%s)\n", err)
		return nil
	}
	log.Debug("Successfully fetched weather forecast")
	return w

}

func init() {
	RegisteredFuncs["weather"] = getWeatherModule
}
