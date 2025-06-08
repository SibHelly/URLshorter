package api

import (
	"net/http"

	"github.com/SibHelly/url_shortnener/internals/app/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(urlhandler *handlers.UrlHandler) *mux.Router {
	r := mux.NewRouter()
	// создать alias.
	r.HandleFunc("/shorten", urlhandler.CreateAlias).Methods("POST")
	// получить информацию про алиас.
	r.HandleFunc("/url/{alias}", func(w http.ResponseWriter, r *http.Request) { return }).Methods("GET")
	// удалить алиас.
	r.HandleFunc("/url/{alias}", func(w http.ResponseWriter, r *http.Request) { return }).Methods("DELETE")
	// перейти по алиасу
	r.HandleFunc("/{alias}", urlhandler.Redirect).Methods("GET")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFound).GetHandler()
	return r
}
