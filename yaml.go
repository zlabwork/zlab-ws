package app

var Yaml = &yaml{}

type yaml struct {
	Access [][]string
	Base   struct {
		MQ     []string `yaml:"mq"`
		LogDir string   `yaml:"logDir"`
	} `yaml:"base"`
}
