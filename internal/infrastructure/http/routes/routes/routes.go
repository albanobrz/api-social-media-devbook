package routes

import (
	"net/http"
)

type Route struct {
	URI        string
	Method     string
	Controller func(http.ResponseWriter, *http.Request)
	NeedAuth   bool
}
