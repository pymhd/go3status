package modules

import (
	"fmt"

	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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
	mo.FullText = fmt.Sprintf("%s%d", mo.FullText, count)	
}

func getDockerCount(v string) (int, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(v))
	if err != nil {
		return 0, err
	}
	defer cli.Close()
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return 0, err
	}
	return len(containers), nil
}

func init() {
	RegisteredFuncs["docker"] = Docker
}
