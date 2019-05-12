package modules

import (
	"fmt"
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"net"
	"strconv"
	"strings"
)

func network(mo *ModuleOutput, cfg ModuleConfig) {
	var conn_level float64
	var addr string

	conn_type := "\uf3dd"
        conn_level = 0

	output := ModuleOutput{}
	output.Markup = "pango"

	active_iface := getActiveIface()
	iface, _ := net.InterfaceByName(active_iface)
	addrs, err := iface.Addrs()
	if err == nil {
		addr = fmt.Sprint(addrs[0])
		if checkWireless(active_iface) {
			wireless := getWirelessInfo(active_iface)
			essid := (wireless["essid"])
			conn_level = getPercent(wireless["level"])
			conn_type = "\uf1eb"
			mo.FullText = fmt.Sprintf("%s %s[%.0f%%] %s", conn_type, essid, conn_level, addr)
		} else {
			conn_type = "\uf0e8"
			conn_level = 100
			mo.FullText = fmt.Sprintf("%s %s", conn_type, addr)
		}
	} else {
		mo.FullText = conn_type
	}

	mo.ShortText = conn_type
	mo.Color = getColor(conn_level, cfg)
}

func getPercent(s string) (percent float64) {
	val := strings.Split(s, "/")
	a, _ := strconv.ParseFloat(val[0], 64)
	b, _ := strconv.ParseFloat(val[1], 64)
	percent = a / b * 100
	return
}

func getActiveIface() (ifname string) {
	re := regexp.MustCompile(`^(\w+)\s+00000000`)

	file, _ := os.Open("/proc/net/route")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		if len(match) != 0 {
			ifname = match[1]
			break
		}

	}
	return
}

func checkWireless(ifname string) (res bool) {
	re := regexp.MustCompile(`^` + ifname)

	file, _ := os.Open("/proc/net/wireless")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		match := re.FindString(scanner.Text())
		if len(match) != 0 {
			res = true
			break
		}
	}
	return
}

func getWirelessInfo(ifname string) (wlan map[string]string) {
        ifcfg := exec.Command("iwconfig", ifname)
        out, _ := ifcfg.CombinedOutput()
	re := regexp.MustCompile(`(?s:ESSID:"(?P<essid>\S+)".*Quality=(?P<level>\S+))`)
        groups := re.SubexpNames()
        match := re.FindStringSubmatch(string(out))
	wlan = make(map[string]string)
        if len(match) == 0 {
                fmt.Println("not match")
                wlan["essid"] = "No ESSID"
        }
        for i, name := range groups {
                if i != 0 && name != "" {
                        wlan[name] = match[i]
                }
        }
	return wlan
}

func init() {
	RegisteredFuncs["network"] = network
}
