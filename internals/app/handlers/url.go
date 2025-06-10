package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/SibHelly/url_shortnener/internals/app/modals"
	urlprocessor "github.com/SibHelly/url_shortnener/internals/app/processors"
	"github.com/gorilla/mux"
)

type UrlHandler struct {
	processor *urlprocessor.UrlProcessor
}

func NewUrlHandler(processor *urlprocessor.UrlProcessor) *UrlHandler {
	handler := new(UrlHandler)
	handler.processor = processor
	return handler
}

func (handler *UrlHandler) CreateAlias(w http.ResponseWriter, r *http.Request) {
	var url modals.Url
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		WrapError(w, err)
		return
	}

	err = handler.processor.CreateAlias(url)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   "",
	}

	WrapOK(w, m)
}

func (hanler *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias := vars["alias"]

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")

	url, err := hanler.processor.GetUrlRedirect(alias)
	if err != nil {
		WrapError(w, err)
		return
	}

	err = hanler.processor.AddVisit(alias)
	if err != nil {
		WrapError(w, err)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
func (hanler *UrlHandler) GetAliases(w http.ResponseWriter, r *http.Request) {
	urls, err := hanler.processor.GetAliases()
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   urls,
	}

	WrapOK(w, m)
}

func (hanler *UrlHandler) GetInfoAboutAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias := vars["alias"]

	url, err := hanler.processor.GetInfoAboutAlias(alias)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   url,
	}

	WrapOK(w, m)

}

func (hanler *UrlHandler) DeleteAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias := vars["alias"]

	err := hanler.processor.DeleteAlias(alias)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   "",
	}

	WrapOK(w, m)

}

// func (hanler *UrlHandler) DeleteAlias(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	alias := vars["alias"]

// 	url, err := hanler.processor.GetInfoAboutAlias(alias)
// 	if err != nil {
// 		WrapError(w, err)
// 		return
// 	}

// 	var m = map[string]interface{}{
// 		"result": "OK",
// 		"data":   url,
// 	}

// 	WrapOK(w, m)

// }
