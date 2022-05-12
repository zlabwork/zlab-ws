package app

var Yaml = &yaml{}

type yaml struct {
	Access [][]string

	Base struct {
		Node    int32  `yaml:"node"`
		Host    string `yaml:"host"`
		Monitor string `yaml:"monitor"`
		PortRpc string `yaml:"portRpc"`
	} `yaml:"base"`

	MQ []string `yaml:"mq"`

	Log string `yaml:"log"`

	Broker struct {
		Port string
	}

	Business struct {
		Cache struct {
			Host string
			Port int64
		}
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
	} `yaml:"business"`
}
