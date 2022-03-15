package restful

import (
	"app"
	"app/msg"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	app.ResponseData(w, msg.OK, app.VersionName+"@"+app.VersionNumber)
}
