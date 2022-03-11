package app

import (
	"app/msg"
	"encoding/json"
	"net/http"
)

type JsonCode struct {
	Code    int
	Message string
}

type JsonData struct {
	Code    int
	Message string
	Data    interface{}
}

func ResponseRaw(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	bs, _ := json.Marshal(data)
	w.Write(bs)
}

func ResponseCode(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	bs, _ := json.Marshal(JsonCode{
		Code:    code,
		Message: msg.Text(code),
	})
	w.Write(bs)
}

func ResponseData(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	bs, _ := json.Marshal(JsonData{
		Code:    code,
		Message: msg.Text(code),
		Data:    data,
	})
	w.Write(bs)
}
