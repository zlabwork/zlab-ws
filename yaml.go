package app

var Yaml = &yaml{}

type yaml struct {
	Access [][]string

	// basic
	Base struct {
		MQ      []string `yaml:"mq"`
		RpcHost string   `yaml:"rpcHost"`
		Monitor string   `yaml:"monitor"`
		LogDir  string   `yaml:"logDir"`
	} `yaml:"base"`

	// cache
	Cache struct {
		Host string
		Port int64
	}

	// database
	Db struct {
		Redis struct {
			Host string
			Port int64
		}
		Mysql struct {
			Host string
			Port int64
			User string
			Pass string
			Name string
		}
	}
}
