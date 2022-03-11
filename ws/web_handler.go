package ws

import (
	"app"
	"app/msg"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	app.ResponseData(w, msg.OK, app.VersionName+"@"+app.VersionNumber)
}

func devHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.EscapedPath() != "/dev" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "../templates/dev.html")
}
