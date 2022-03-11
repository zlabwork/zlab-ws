package msg

const (
	OK           = 200
	Err          = 400
	ErrNotFound  = 404
	ErrTimeout   = 20001
	ErrSignature = 20002
	ErrAccess    = 20003
	ErrEncode    = 20004
	ErrHeader    = 20005
	ErrParameter = 20006
	ErrProcess   = 20007
	ErrNoData    = 20008
	ErrModify    = 20009
)

var statusText = map[int]string{
	OK:           "success",
	Err:          "error",
	ErrNotFound:  "page not found",
	ErrTimeout:   "error request time",
	ErrSignature: "error request signature",
	ErrAccess:    "error access",
	ErrEncode:    "error encode",
	ErrHeader:    "error request headers",
	ErrParameter: "error parameter",
	ErrProcess:   "error in execute process",
	ErrNoData:    "can not find",
	ErrModify:    "error when modify data",
}

func Text(code int) string {
	return statusText[code]
}
