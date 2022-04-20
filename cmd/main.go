package main

import (
	"app"
	"app/service"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"time"
)

func init() {

	// env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	// app.yaml
	bs, err := ioutil.ReadFile("../config/app.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if yaml.Unmarshal(bs, app.Yaml) != nil {
		log.Fatal(err)
	}

	// logs
	_ = os.MkdirAll(app.Yaml.Base.LogDir, 0666)
	f, err := os.OpenFile(app.Yaml.Base.LogDir+string(os.PathSeparator)+"logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(f)
	}
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})

	// libs
	app.Libs = app.NewLibs()
}

func main() {

	// params
	var wait time.Duration
	var module string
	var addr = flag.String("addr", ":8080", "http service address")
	flag.StringVar(&module, "module", "", "module name - e.g. broker or business")
	flag.DurationVar(&wait, "timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	if len(os.Getenv("APP_PORT")) > 0 {
		*addr = ":" + os.Getenv("APP_PORT")
	}

	// service
	switch module {
	case "business":
		srv, err := service.NewBusinessService()
		if err != nil {
			log.Fatal(err)
		}
		srv.Run()

	case "broker":
		srv, err := service.NewBrokerService()
		if err != nil {
			log.Fatal(err)
		}
		srv.Run(addr)

	default:
		srv, err := service.NewBrokerService()
		if err != nil {
			log.Fatal(err)
		}
		srv.Run(addr)
	}
	app.Banner("service is started")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	fmt.Println("logs at " + app.Yaml.Base.LogDir)
	log.Println("shutting down")
	os.Exit(0)
}
