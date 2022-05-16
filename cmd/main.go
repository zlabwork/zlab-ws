package main

import (
	"app"
	"app/restful"
	"app/service"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func init() {

	// env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(err)
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
	_ = os.MkdirAll(app.Yaml.Log, 0666)
	f, err := os.OpenFile(app.Yaml.Log+string(os.PathSeparator)+"logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(f)
	}
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})

	// libs
	app.Libs = app.NewLibs(app.Yaml.Base.Node)
}

func usage() {

	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

func login() {

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", restful.DefaultHandler)
	r.HandleFunc("/auth/{token:[0-9a-zA-Z_-]+}", restful.AuthHandler).Methods(http.MethodGet)

	addr := "0.0.0.0:" + strconv.FormatInt(app.Yaml.Login.Port, 10)
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	fmt.Println("login service at " + addr)
	log.Println("login service at " + addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

func main() {

	// params
	var help bool
	var wait time.Duration
	var module string
	var addr = flag.String("addr", ":"+app.Yaml.Broker.Port, "http service address")
	flag.Usage = usage
	flag.BoolVar(&help, "h", false, "help")
	flag.StringVar(&module, "m", "", "module name - e.g. -m broker, -m business or -m dev")
	flag.DurationVar(&wait, "timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	if help {
		usage()
		os.Exit(1)
	}

	// banner
	app.Banner("service is started")

	// service
	switch module {
	case "login":
		login()

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

	case "monitor":
		srv, err := service.NewMonitorService()
		if err != nil {
			log.Fatal(err)
		}
		srv.Run()

	case "dev":
		login()

		mon, err := service.NewMonitorService()
		if err != nil {
			log.Fatal(err)
		}
		mon.Run()

		bus, err := service.NewBusinessService()
		if err != nil {
			log.Fatal(err)
		}
		bus.Run()

		bro, err := service.NewBrokerService()
		if err != nil {
			log.Fatal(err)
		}
		bro.Run(addr)

	default:
		usage()
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	fmt.Println("logs at " + app.Yaml.Log)
	log.Println("shutting down")
	os.Exit(0)
}
