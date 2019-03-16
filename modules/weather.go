package modules


import (
        "fmt"
        "time"
        "strings"
        "strconv"
        "net/http"
        "sync/atomic"
        "encoding/json"
)
const (
    OpenWeatherMapToken = "af7bfe22287c652d032a3064ffa44088"
)

type location struct {
    Name	string	`json:"city"`
    Coord	string  `json:"loc"`
    rewrite	string
}

type weather struct {
    Weather  []struct {
            Main  string `json:"main"`
            Desc  string  `json:"description"`    
    } `json:"weather"`
    Main struct {
        Temp float64 `json:"temp"`
    } `json:"main"`	
    Wind  struct {
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

func(w Weather) run(c chan ModuleOutput, cfg ModuleConfig) {
        output := ModuleOutput{}
	output.Name = w.name
	output.Instance = strconv.Itoa(cfg.Id)
	output.Refresh = true
	output.Markup = "pango"
	output.FullText = cfg.Prefix
	
	loc := new(location)
	l, ok := cfg.Extra["location"]
	if !ok {
	    loc = getLocation()
	} else {
	    city, _ := l.(string)
	    if city == "auto" {
	        loc = getLocation()
	    } else {
	        //loc = &location{}
	        loc.rewrite = city 
	        loc.Name = city
            }
	} 
	wf := getWeather(loc)
	forecast := fmt.Sprintf("%s %s %f %d", loc.Name, wf.Weather[0].Desc, wf.Main.Temp, wf.Wind.Speed)
	
	if x := atomic.LoadInt32(Mute[cfg.Id]); x == -1 {
                output.FullText += "..."
        } else {
                output.FullText += forecast
        }

	c <- output
}

func (w Weather) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {

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
            fmt.Println("Not worked out")
        }
        return loc
}

func getWeather(l *location) *weather{
        w := new(weather)
        coord := strings.Split(l.Coord, ",")
        url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&lat=%s&lon=%s&units=metric", OpenWeatherMapToken, coord[0], coord[1])
        res, err := http.Get(url)
        if err != nil {
            return w
        }
        defer res.Body.Close()
        if err := json.NewDecoder(res.Body).Decode(w); err != nil {
            fmt.Println("Not worked out")
        }
        return w


}

func init() {
    w := Weather{name: "weather"}
    selfRegister(w)
}