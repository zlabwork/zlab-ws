package service

import "app/service/control"

type Monitor struct {
}

func NewMonitorService() (*Monitor, error) {

	return &Monitor{
	}, nil
}

func (mo *Monitor) Run() {

	go func() {
		control.RunMonitorServer()
	}()
}
