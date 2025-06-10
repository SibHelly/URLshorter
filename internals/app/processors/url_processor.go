package urlprocessor

import (
	"errors"
	"fmt"
	"time"

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
	if alias.Visit_count <= 0 {
		return errors.New("visit count should not be negative or equal 0")
	}
	if alias.Expires_at != nil {
		if alias.Expires_at.Before(time.Now()) {
			return errors.New("expires time should not be before time now")
		}
	}

	return processor.storage.CreateAlias(alias)
}

func (processor *UrlProcessor) DeleteAlias(alias string) error {
	if alias == "" {
		return errors.New("alias should not be empty")
	}

	return processor.storage.DeleteAlias(alias)
}

func (processor *UrlProcessor) GetUrlRedirect(alias string) (string, error) {
	if alias == "" {
		return "", errors.New("alias should not be empty")
	}
	url, err := processor.storage.GetUrl(alias)
	if err != nil {
		return "", err
	}
	err = processor.storage.CheckVisitCount(alias)
	if err != nil {
		return "", err
	}
	err = processor.storage.CheckExpiresAt(alias)
	if err != nil {
		return "", err
	}
	err = processor.storage.UpdateVisitCount(alias)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (processor *UrlProcessor) GetInfoAboutAlias(alias string) (*modals.UrlInfo, error) {
	if alias == "" {
		return nil, errors.New("alias should not be empty")
	}
	return processor.storage.GetInfoAboutAlias(alias)
}

func (processor *UrlProcessor) GetAliases() ([]modals.Url, error) {
	return processor.storage.GetAliases()
}

func (processor *UrlProcessor) AddVisit(alias string) error {
	if alias == "" {
		return errors.New("alias should not be empty")
	}
	url_id, err := processor.storage.GetUrlId(alias)
	if err != nil {
		return fmt.Errorf("Can't get id this alias err: %w", err)
	}
	err = processor.storage.AddVisit(url_id)
	if err != nil {
		return fmt.Errorf("Can't add visit to alias err: %w", err)
	}
	return nil
}
