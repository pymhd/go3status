package modules

func cpu(mo *ModuleOutput, cfg ModuleConfig) {
        mo.FullText = "70%"
}

func init() {
        RegisteredFuncs["cpu"] = cpu
}
