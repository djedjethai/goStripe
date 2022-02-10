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

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "userID") {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
		next.ServeHTTP(w, r)
	})
}
