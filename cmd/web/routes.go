package main

import (
	"net/http"

	"github.com/justinas/alice"
	"textvault/ui"
)

func (app *application) routers() http.Handler {
	mux := http.NewServeMux()
	//добавили получение сторонних файлов из static
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	//добавили цепочку session
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticated)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc((app.snippetCreate)))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /snippet/delete/{id}", protected.ThenFunc(app.deletePost))

	// добавляем новые роуты
	mux.Handle("GET /user/signup/", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup/", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standart := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standart.Then(mux)

}
