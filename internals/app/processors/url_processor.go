package urlprocessor

import (
	"errors"
	"fmt"

	"github.com/SibHelly/url_shortnener/internals/app/db"
	"github.com/SibHelly/url_shortnener/internals/app/modals"
)

type UrlProcessor struct {
	storage *db.UrlStorage
}

func NewUrlProcessor(storage *db.UrlStorage) *UrlProcessor {
	processor := new(UrlProcessor)
	processor.storage = storage
	return processor
}

func (processor *UrlProcessor) CreateAlias(alias modals.Url) error {
	if alias.Original_url == "" {
		return errors.New("url should not be empty")
	}
	if alias.Alias == "" {
		return errors.New("alias should not be empty")
	}

	fmt.Println(alias)

	return processor.storage.CreateAlias(alias)
}

func (processor *UrlProcessor) GetUrl(alias string) (string, error) {
	if alias == "" {
		return "", errors.New("alias should not be empty")
	}
	return processor.storage.GetUrl(alias)
}
