package restful

import (
	"app"
	"app/msg"
	"app/service"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token := vars["token"]
	fmt.Println(token)

	srv, err := service.NewLoginService()
	if err != nil {

	}

	if !srv.CheckToken(token) {
		app.ResponseCode(w, msg.Err)
		return
	}

	app.ResponseData(w, msg.OK, "success")
}
