package modules

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func memory(mo *ModuleOutput, cfg ModuleConfig) {
	var total, avail string
	var done int
	meminfo, _ := os.Open("/proc/meminfo")
	defer meminfo.Close()

	scanner := bufio.NewScanner(meminfo)
	for scanner.Scan() {
		line := scanner.Text()
		sl := strings.Fields(line)
		switch sl[0] {
		case "MemTotal:":
			total = sl[1]
			done += 1
		case "MemAvailable:":
			avail = sl[1]
			done += 1
		}
		if done == 2 {
			break
		}
	}
	totalF, _ := strconv.ParseFloat(total, 64)
	availF, _ := strconv.ParseFloat(avail, 64)
	usedF := totalF - availF
	percentage := 100 * (usedF / totalF)

	mo.FullText = fmt.Sprintf("%s%.1f/%.1f (%.0f%%)", mo.FullText, usedF/1048576, totalF/1048576, percentage)
}

func init() {
	RegisteredFuncs["memory"] = memory
}
