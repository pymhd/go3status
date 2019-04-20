package modules

import (
        "fmt"
        "net"
        "net/http"
        "io/ioutil"
        "encoding/json"
)


func Docker(mo *ModuleOutput, cfg ModuleConfig) {
	v, ok := cfg.Extra["clientAPIVersion"]
	if !ok {
		mo.FullText += "Unknown version"
		return
	}
	ver, ok := v.(string)
	if !ok {
		mo.FullText += "Version must be string"
		return
	}
	cv, ok := cfg.Extra["color"]
	if ok {
		mo.Color, _ = cv.(string)
	}
	count, err := getDockerCount(ver)
	if err != nil {
		mo.FullText += "Daemon OFF"
		return
	}
	mo.ShortText = cfg.Prefix
	mo.FullText = fmt.Sprintf("%s%d", mo.FullText, count)
}


func getDockerCount(version string) (int, error){
        var tr http.Transport
        v := cache.Get("docker:transport")
        if v == nil {
                tr = http.Transport{ Dial: func(string, string) (net.Conn, error) {
                        return net.Dial("unix", "/var/run/docker.sock")
                }}
                cache.Add("docker:transport", tr, "1h")
        } else {
                tr, _ = v.(http.Transport)
        }
        defer tr.CloseIdleConnections()
        client := &http.Client{Transport: &tr}
        url := fmt.Sprintf("http://%s/containers/json", version)
        resp, err := client.Get(url)
        if err != nil {
                return 0, err
        }
        defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        var d []interface{}
        err = json.Unmarshal(body, &d)
        if err != nil {
                return 0, err
        } 
        return len(d), nil
}    

func init() {
        RegisteredFuncs["docker"] = Docker
}
