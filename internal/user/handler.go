package user

import "net/http"

type HandlerUser struct {
	*HandlerUserDep
}
type HandlerUserDep struct {
	*ServiceUser
}

func NewHandlerUser(router *http.ServeMux, dep *HandlerUserDep) {

}
