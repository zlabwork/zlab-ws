package main

import (
	"io/ioutil"
	"log"
	"zlabws"
	"zlabws/ws"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	bs, err := ioutil.ReadFile("../configs/app.yaml")
	err = yaml.Unmarshal(bs, &zlabws.Cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	ws.Run()
}
