package modules

import (
        "strconv"
)

func cpu(mo *ModuleOutput, cfg ModuleConfig) {
        var percentage int
        v := cache.Get("cpu:percentage")
        if v != nil {
                percentage, _ = v.(int)
        }
        mo.FullText += strconv.Itoa(percentage)
        percentage++
        cache.Add("cpu:percentage", percentage, "1h")
}

func init() {
        RegisteredFuncs["cpu"] = cpu
}
