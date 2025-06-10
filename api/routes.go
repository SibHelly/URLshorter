package api

import (
	"github.com/SibHelly/url_shortnener/internals/app/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(urlhandler *handlers.UrlHandler) *mux.Router {
	r := mux.NewRouter()
	// создать alias.
	r.HandleFunc("/shorten", urlhandler.CreateAlias).Methods("POST")
	// получить информацию про алиас.
	r.HandleFunc("/url/{alias}", urlhandler.GetInfoAboutAlias).Methods("GET")
	// удалить алиас.
	r.HandleFunc("/url/{alias}", urlhandler.DeleteAlias).Methods("DELETE")
	// перейти по алиасу
	r.HandleFunc("/{alias}", urlhandler.Redirect).Methods("GET")
	// получить все алиасы
	r.HandleFunc("/my/all", urlhandler.GetAliases).Methods("GET")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFound).GetHandler()
	return r
}
