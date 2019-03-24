package modules

import (
	"fmt"
	"os"
        "time"
	"bufio"
	"regexp"
	"strconv"
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

	percentage, capacity, status := getBatPercent()
        for lvl, val := range cfg.Levels {
		if inRange(percentage, val) {
			output.Color = cfg.Colors[lvl]
                }
	}
        duration, _ := time.ParseDuration(capacity)
        hours, minutes := fmtDuration(duration)
	output.FullText += fmt.Sprintf("%.0f%% %s(%dh%dm)", percentage, status, hours, minutes)
        c <- output
}

func (batt BAT) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
}

func fmtDuration(d time.Duration) (time.Duration, time.Duration) {
    d = d.Round(time.Minute)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    return h, m /*fmt.Sprintf("%02dh %02dm", h, m) */
}


func getBatPercent () (float64, string, string) {
        var b_status, capacity string
	var e_full, e_now, p_now int
	var percent float64

	re := regexp.MustCompile(`(^.*)=(.*$)`)
	file, _ := os.Open("/sys/class/power_supply/BAT0/uevent")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if match != nil {
			switch match[1] {
			case "POWER_SUPPLY_STATUS":
				b_status = match[2]
			case "POWER_SUPPLY_POWER_NOW":
				p_now, _ = strconv.Atoi(match[2])
			case "POWER_SUPPLY_ENERGY_NOW":
				e_now, _ = strconv.Atoi(match[2])
			case "POWER_SUPPLY_ENERGY_FULL":
				e_full, _ = strconv.Atoi(match[2])
			}
		}
	}

	percent = float64(100 * e_now / e_full)
	switch b_status {
	case "Discharging":
		c := float64(e_now) / float64(p_now)
		capacity = fmt.Sprintf("%.7fh", c)
		b_status = "DIS"
	case "Charging":
		c := (float64(e_full) - float64(e_now)) / float64(p_now)
		capacity = fmt.Sprintf("%.7fh", c)
		b_status = "CHR"
	}

	return percent, capacity, b_status
}

func init() {
	batt := BAT{name: "batt"}
	selfRegister(batt)
}
