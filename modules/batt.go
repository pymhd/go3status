package modules

import (
	"fmt"
	"os"
        "time"
)

type BAT struct {
	name string
}

func (batt BAT) Name() string {
        return batt.name
}

func (batt BAT) Run(c chan ModuleOutput, cfg ModuleConfig) {
        //to run by start
        batt.run(c, cfg)

        // to run periodically
        ticker := time.NewTicker(cfg.Interval)
        for {
                select {
                case <-ticker.C:
                        batt.run(c, cfg)
                case <-RefreshChans[cfg.Id]:
                        batt.run(c, cfg)
                }
        }
}

func (batt BAT) run(c chan ModuleOutput, cfg ModuleConfig) {
        output := ModuleOutput{}
	output.Name = batt.name
        output.Markup = "pango"
        output.Refresh = true
        output.FullText = cfg.Prefix

	percentage, capacity := getBatPercent()
        for lvl, val := range cfg.Levels {
		if inRange(percentage, val) {
			output.Color = cfg.Colors[lvl]
                }
	}
        duration, _ := time.ParseDuration(capacity)
        tu := fmtDuration(duration)
        output.FullText += fmt.Sprintf("%.2f%s %s", percentage, cfg.Postfix, tu)
        c <- output
}

func (batt BAT) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
}

func fmtDuration(d time.Duration) string {
    d = d.Round(time.Minute)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    return fmt.Sprintf("%02dh %02dm", h, m)
}


func getBatPercent () (float64, string) {
	var full, now, power int

	data, _ := os.Open("/sys/class/power_supply/BAT0/energy_full")
	fmt.Fscanf(data, "%d", &full)

	data, _ = os.Open("/sys/class/power_supply/BAT0/energy_now")
	fmt.Fscanf(data, "%d", &now)

        data, _ = os.Open("/sys/class/power_supply/BAT0/power_now")
        fmt.Fscanf(data, "%d", &power)

	res := 100 * now / full
        percent := float64(res)

        resf := float64(now) / float64(power)
        capacity := fmt.Sprintf("%.4fh", resf)
	return percent, capacity
}

func init() {
	batt := BAT{name: "batt"}
	selfRegister(batt)
}
