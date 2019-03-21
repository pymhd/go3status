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
        tu, _ := time.ParseDuration(capacity)
        output.FullText += fmt.Sprintf("%.2f%s %.0fh:%.0fm", percentage, cfg.Postfix, tu.Hours(), tu.Minutes())
        c <- output
}

func (batt BAT) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
}

func getBatPercent () (float64, string) {
	var full, now, power int

	data, _ := os.Open("/sys/class/power_supply/BAT0/energy_full")
	fmt.Fscanf(data, "%d", &full)

	data, _ = os.Open("/sys/class/power_supply/BAT0/energy_now")
	fmt.Fscanf(data, "%d", &now)

        data, _ = os.Open("/sys/class/power_supply/BAT0/power_now")
        fmt.Fscanf(data, "%d", &power)
        power = 12312123

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
