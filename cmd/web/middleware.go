package main

import (
	"net/http"
)

// this is a middleware from the package alexedwards/scs/v2
// receive an httpHandler, we modify it, and the return an httpHandler
// session is a package level variable(so we have access to it)
// then add the middleware in the router(mux)
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
