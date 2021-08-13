package main

import (
    "github.com/joho/godotenv"
    "log"
    "zlabws/ws"
)

func main() {

    err := godotenv.Load("../.env")
    if err != nil {
        log.Fatal(err)
    }

    ws.Run()
}
