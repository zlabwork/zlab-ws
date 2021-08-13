package main

import (
    "github.com/joho/godotenv"
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "log"
    "zlabws"
    "zlabws/ws"
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
