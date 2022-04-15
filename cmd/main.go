package main

import (
	"app"
	"app/restful"
	"app/service"
	"app/service/broker"
	"app/service/business"
	"flag"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

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

	// libs
	app.Libs = app.NewLibs()

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

	case "broker":
		broker.Main()
		log.Println("Broker is started")

	case "business":
		business.Main()
		log.Println("Business is started")

	default:
		srv, err := service.NewService()
		if err != nil {
			log.Fatal(err)
		}
		srv.Run()

		http.HandleFunc("/", restful.DefaultHandler)
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			srv.ServeWs(w, r)
		})

		// Run our server in a goroutine so that it doesn't block.
		go func() {
			if err := http.ListenAndServe(*addr, nil); err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		}()
		app.Banner("Service port :" + os.Getenv("APP_PORT"))
		log.Println("Service is started")
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
	log.Println("shutting down")
	os.Exit(0)

}
