package modules

import (
        "os/exec"
        "strings"
        "strconv"
)

var (
        clickMap =  map[int]string{1: "left", 3: "right", 4: "wheelUp", 5: "wheelDown"} 
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


func execute(cmd string) {
        if len(cmd) > 0 {
                args := strings.Split(cmd, " ")
                c := exec.Command(args[0], args[1:]...)
                c.Start()
        }
}

