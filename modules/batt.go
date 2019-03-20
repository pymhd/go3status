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

	percentage := getBatPercent()
        for lvl, val := range cfg.Levels {
		if inRange(percentage, val) {
			output.Color = cfg.Colors[lvl]
                }
	}
        output.FullText += fmt.Sprintf("%.2f%s", percentage, cfg.Postfix)
        c <- output
}

func (batt BAT) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
}

func getBatPercent () float64 {
	var full, now int

	data, _ := os.Open("/sys/class/power_supply/BAT0/energy_full")
	fmt.Fscanf(data, "%d", &full)

	data, _ = os.Open("/sys/class/power_supply/BAT0/energy_now")
	fmt.Fscanf(data, "%d", &now)

	res := 100 * now / full
	return float64(res)
}

func init() {
	batt := BAT{name: "batt"}
	selfRegister(batt)
}
