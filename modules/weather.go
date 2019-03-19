package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
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

	loc *location
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

func (w Weather) Name() string {
	return w.name
}

func (w Weather) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//to run by start
	w.run(c, cfg)

	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			w.run(c, cfg)
		case <-RefreshChans[cfg.Id]:
			w.run(c, cfg)
		}
	}
}

func (w Weather) run(c chan ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}
	output.Name = w.name
	output.Instance = strconv.Itoa(cfg.Id)
	output.Refresh = true
	output.Markup = "pango"
	output.Separator = true
	output.FullText = cfg.Prefix

	//if it is first time
	if loc == nil {
		l, ok := cfg.Extra["location"]
		if !ok {
			loc = getLocation()
		} else {
			city, _ := l.(string)
			if city == "auto" {
				loc = getLocation()
			} else {
				loc = new(location)
				loc.rewrite = city
				loc.Name = city
			}
		}
	}
	wf := getWeather(loc)
	icon, ok := iconSet[wf.Weather[0].Main]
	if !ok {
		icon = smogIcon
	}
	forecast := fmt.Sprintf("%s: %s %.0f%s (%d m/s)", loc.Name, icon, wf.Main.Temp, celsiusIcon, wf.Wind.Speed)

	if x := atomic.LoadInt32(Mute[cfg.Id]); x == -1 {
		output.FullText += forecast[strings.Index(forecast, ":")+2:]
	} else {
		output.FullText += forecast
	}

	c <- output
}

func (w Weather) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
	switch ce.Button {
	// middle, reserved, shrink panel and force refresh
	case 2:
		w.Mute(cfg.Id)
		RefreshChans[cfg.Id] <- true
	// any other
	default:
		buttonNumber := ce.Button
		buttonText := clickMap[buttonNumber]
		cmd, ok := cfg.ClickEvents[buttonText]
		if !ok {
			//if no cmd specified in config file
			break
		}
		execute(cmd)
		RefreshChans[cfg.Id] <- true

	}
}

func (w Weather) Mute(id int) {
	atomic.StoreInt32(Mute[id], ^atomic.LoadInt32(Mute[id]))
}

func getLocation() *location {
	loc := new(location)
	url := "http://ipinfo.io/json"
	res, err := http.Get(url)
	if err != nil {
		return loc
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(loc); err != nil {
		panic(err)
	}
	return loc
}

func getWeather(l *location) *weather {
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
	w := Weather{name: "weather"}
	selfRegister(w)
}
