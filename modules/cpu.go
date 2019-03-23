package modules

import (
	"fmt"
	"os"
)

const (
	cpuTimeKey = "cpu:time"
)

type cpuTime struct {
	user, nice, system, idle, io, irq, softirq, steal, guest, guest_nice int
}

func cpu(mo *ModuleOutput, cfg ModuleConfig) {
	stat, _ := os.Open("/proc/stat")
	defer stat.Close()

	prevCpuTime := cpuTime{}
	newCpuTime := cpuTime{}
	var tmp string

	//Get prev cpu times
	cv := cache.Get(cpuTimeKey)
	if cv != nil {
		//value exist
		prevCpuTime, _ = cv.(cpuTime)
	}

	//get new cpu times
	fmt.Fscanf(stat, "%s %d %d %d %d %d %d %d %d %d %d", &tmp, &newCpuTime.user, &newCpuTime.nice, &newCpuTime.system, &newCpuTime.idle, &newCpuTime.io, &newCpuTime.irq, &newCpuTime.softirq, &newCpuTime.steal, &newCpuTime.guest, &newCpuTime.guest_nice)

	PrevIdle := prevCpuTime.idle + prevCpuTime.io
	Idle := newCpuTime.idle + newCpuTime.io

	PrevNonIdle := prevCpuTime.user + prevCpuTime.nice + prevCpuTime.system + prevCpuTime.irq + prevCpuTime.softirq + prevCpuTime.steal
	NonIdle := newCpuTime.user + newCpuTime.nice + newCpuTime.system + newCpuTime.irq + newCpuTime.softirq + newCpuTime.steal

	PrevTotal := PrevIdle + PrevNonIdle
	Total := Idle + NonIdle

	totald := Total - PrevTotal
	idled := Idle - PrevIdle

	cpu := 100 * float64(totald-idled) / float64(totald)

	cache.Add(cpuTimeKey, newCpuTime, "1h")
	//Generate output
	mo.Color = getColor(cpu, cfg)
	mo.FullText = fmt.Sprintf(" %.2f%%", cpu)
}


func init() {
	RegisteredFuncs["cpu"] = cpu
}
