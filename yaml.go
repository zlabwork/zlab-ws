package app

var Yaml = &yaml{}

type yaml struct {
	Access [][]string
	Base   struct {
		LogDir string `yaml:"logDir"`
	} `yaml:"base"`
}
