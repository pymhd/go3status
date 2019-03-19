package batt

import (
	"fmt"
	"os"
)

type BAT struct {
	name string
}

func (batt BAT) run() {
	output.Name = batt.name

	persentage := getBatPercent(batt.name)
	
}

func getBatPercent (Name) int {
	var full, now, res int

	data, _ := os.Open("/sys/class/power_supply/" + Name + "/energy_full")
	fmt.Fscanf(data, "%d", &full)

	data, _ = os.Open("/sys/class/power_supply/" + Name + "/energy_now")
	fmt.Fscanf(data, "%d", &now)

	res = 100 * now / full
	return res
}

func init() {
	batt := BAT{name: "BAT0"}
	selfRegister(batt)
}
