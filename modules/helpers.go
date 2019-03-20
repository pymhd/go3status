package modules

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	clickMap = map[int]string{1: "left", 3: "right", 4: "wheelUp", 5: "wheelDown"}
)

func inRange(p float64, r string) bool {
	vals := strings.Split(r, "-")
	if len(vals) != 2 {
		return false
	}
	min, err := strconv.ParseFloat(vals[0], 64)
	if err != nil {
		return false
	}
	max, err := strconv.ParseFloat(vals[1], 64)
	if err != nil {
		return false
	}
	if p >= min && p <= max {
		return true
	}
	return false

}

func execute(oneliner string, timeout time.Duration) string {
	if len(oneliner) == 0 {
		return "Wrong cmd"
	}
	res := make(chan string)
	go func(c chan string) {
		out, err := exec.Command("bash", "-c", oneliner).Output()
		if err != nil {
			c <- fmt.Sprintf("Failed to exec (%s) (%s)", oneliner, err)
			return
		}
		c <- fmt.Sprintf("%s", bytes.Trim(out, " \n,:\t\""))
	}(res)

	select {
	case result := <-res:
		return result
	case <-time.After(timeout):
		return ""
	}
}
