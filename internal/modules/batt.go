package modules

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

func batt(mo *ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}
	output.Markup = "pango"
	//output.Refresh = true
	output.FullText = cfg.Prefix

	percentage, capacity, status := getBatPercent()
	duration, _ := time.ParseDuration(capacity)
	hours, minutes := fmtDuration(duration)
	mo.Color = getColor(percentage, cfg)
	mo.FullText += fmt.Sprintf("%s %.0f%%/%dh%dm", status, percentage, hours, minutes)
	mo.ShortText = fmt.Sprintf("%s", status)
}

func fmtDuration(d time.Duration) (time.Duration, time.Duration) {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return h, m /*fmt.Sprintf("%02dh %02dm", h, m) */
}

func getBatPercent() (float64, string, string) {
	var b_status, capacity string
	var e_full, e_now, p_now, percent float64

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
				p_now, _ = strconv.ParseFloat(match[2], 64)
			case "POWER_SUPPLY_VOLTAGE_NOW":
				p_now, _ = strconv.ParseFloat(match[2], 64)
			case "POWER_SUPPLY_ENERGY_NOW":
				e_now, _ = strconv.ParseFloat(match[2], 64)
			case "POWER_SUPPLY_ENERGY_FULL":
				e_full, _ = strconv.ParseFloat(match[2], 64)
			}
		}
	}

	percent = 100 * e_now / e_full
	switch b_status {
	case "Discharging":
		c := e_now / p_now
		capacity = fmt.Sprintf("%.7fh", c)
		b_status = "\uf242"
	case "Charging":
		c := (e_full - e_now) / p_now
		capacity = fmt.Sprintf("%.7fh", c)
		b_status = "\uf1e6"
	case "Full":
		c := (e_full - e_now) / p_now
		capacity = fmt.Sprintf("%.7fh", c)
		b_status = "\uf240"
	case "Unknown":
		b_status = "\uf5df"
	}

	return percent, capacity, b_status
}

func init() {
	RegisteredFuncs["batt"] = batt
}
