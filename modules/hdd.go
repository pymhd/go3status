package modules

import (
	"bytes"
	"syscall"
	"text/template"
)

const (
	BaseTemplate   = `{{printf "%.1f" $used}}/{{printf "%.1f" $total}} (printf "%.0f" $percentage)`
	TemplatePrefix = `{{$path := .Path}}{{$used := .Used}}{{$total := .Total}}{{$avail := .Avail}}{{$percentage := .Percentage}}`
)

type filesystem struct {
	Path       string
	Used       float64
	Avail      float64
	Total      float64
	Percentage float64
}

func hdd(mo *ModuleOutput, cfg ModuleConfig) {
	fs := filesystem{}
	mp, ok := cfg.Extra["mountPoint"]
	if !ok {
		mo.FullText += "N/A"
		return
	}
	fs.Path, _ = mp.(string)
	fs.Used, fs.Total = getMpStats(fs.Path)
	fs.Percentage = 100 * fs.Used / fs.Total
	fs.Avail = fs.Total - fs.Used
	//generating fulltext
	var T string
	tpl, ok := cfg.Extra["format"]
	if !ok {
		T = BaseTemplate
	} else {
		T, _ = tpl.(string)
		if len(T) == 0 {
			//empty format key
			T = BaseTemplate
		}
	}
	//to register vars
	T = TemplatePrefix + T
	t := template.Must(template.New(fs.Path).Parse(T))

	var out bytes.Buffer
	t.Execute(&out, fs)

	mo.Color = getColor(fs.Percentage, cfg)
	if cfg.ShortFormat {
		mo.FullText = cfg.Prefix
	} else {
		mo.FullText += out.String() + " Gb"
	}
}

func getMpStats(path string) (float64, float64) {
	stats := syscall.Statfs_t{}
	syscall.Statfs(path, &stats)
	avail := float64(stats.Bavail) * float64(stats.Bsize)
	total := float64(stats.Blocks) * float64(stats.Bsize)
	return (total - avail) / 1024 / 1024 / 1024, total / 1024 / 1024 / 1024
}

func init() {
	RegisteredFuncs["hdd"] = hdd
}
