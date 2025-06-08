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

	url, err := hanler.processor.GetUrl(alias)
	if err != nil {
		WrapError(w, err)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
