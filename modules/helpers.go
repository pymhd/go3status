package modules

import (
        "strings"
        "strconv"
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
