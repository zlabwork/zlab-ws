package app

var Yaml = &yaml{}

type yaml struct {
	Access [][]string
	Base   struct {
		MQ      []string `yaml:"mq"`
		RpcHost string   `yaml:"rpcHost"`
		LogDir  string   `yaml:"logDir"`
	} `yaml:"base"`
}
