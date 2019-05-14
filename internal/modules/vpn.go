package modules

import (
	"fmt"
	"regexp"
	"net"
)

func vpn(mo *ModuleOutput, cfg ModuleConfig) {
	var addrs []string
	var vpn_status float64

	output := ModuleOutput{}
	output.Markup = "pango"
	mo.ShortText = "\uf57c"

	re := regexp.MustCompile(`(vpn|tun)`)

	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		if re.MatchString(i.Name) {
			ip, err :=  i.Addrs()
			if err == nil {
				addrs = append(addrs, fmt.Sprint(ip[0]))
			}
		}
	}
	if len(addrs) != 0 {
		vpn_status = 100
                mo.ShortText = "\uf57d"
	}
	mo.Color = getColor(vpn_status, cfg)
	mo.FullText = fmt.Sprintf("%s %s", mo.ShortText, addrs) 
}

func init() {
	RegisteredFuncs["vpn"] = vpn
}
