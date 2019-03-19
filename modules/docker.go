package modules

import (
        "fmt"
        "time"
        "strconv"

        "context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

)

type DockerClient struct {
        name	string
}

func (dc DockerClient) Name() string {
        return dc.name
}

func (dc DockerClient) Run(c chan ModuleOutput, cfg ModuleConfig) {
        //to run by start
        dc.run(c, cfg)

        // to run periodically
        ticker := time.NewTicker(cfg.Interval)
        for {
                select {
                case <-ticker.C:
                        dc.run(c, cfg)
                case <-RefreshChans[cfg.Id]:
                        dc.run(c, cfg)
                }
        }
}

func (dc DockerClient) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {

}

func (dc DockerClient) run(c chan ModuleOutput, cfg ModuleConfig) {
        output := ModuleOutput{}
        output.Name = dc.name
        output.Instance = strconv.Itoa(cfg.Id)
        output.Refresh = true
        output.Markup = "pango"
        output.FullText = cfg.Prefix

        v, ok := cfg.Extra["clientAPIVersion"]
        if !ok {
                output.FullText += "Unknown version"
                c <- output
                return
        }
        version, ok  := v.(string)
        if !ok {
                output.FullText += "Unknown version fmt"
                c <- output
                return
        }
        count, err := getDockerCount(version)
        if err != nil {
                output.FullText += "Daemon OFF"
        } else {
                output.FullText = fmt.Sprintf("%s %d", output.FullText, count) 
        }
        color, ok := cfg.Extra["color"]
        if ok {
                output.Color, _ = color.(string)
        }
        c <- output
}

func getDockerCount(v string) (int, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(v))
	if err != nil {
		return 0, err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return 0, err
	}
	return len(containers), nil
	//for _, container := range containers {
	//	fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	//}
}


func init() {
        dc := DockerClient{name: "docker"}
        selfRegister(dc)
}

