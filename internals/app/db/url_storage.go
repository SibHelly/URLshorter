package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SibHelly/url_shortnener/internals/app/modals"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlStorage struct {
	databasePool *pgxpool.Pool
}

func NewUrlStorage(pool *pgxpool.Pool) *UrlStorage {
	storage := new(UrlStorage)
	storage.databasePool = pool
	return storage
}

func (storage *UrlStorage) GetUrl(alias string) (string, error) {
	var url string
	query := "SELECT original_url FROM urls WHERE alias = $1"
	err := storage.databasePool.QueryRow(context.Background(), query, alias).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("Url for this alias not found\nerr:%w", err)
	}
	return url, nil
}

func (storage *UrlStorage) CreateAlias(alias modals.Url) error {
	var columns []string
	var placeholders []string
	var values []interface{}
	paramIndex := 1

	// Карта всех возможных полей с их значениями и условиями
	fieldsMap := map[string]interface{}{
		"original_url": alias.Original_url, // всегда добавляется
		"alias":        alias.Alias,        // всегда добавляется
		"created_at":   time.Now(),         // всегда добавляется
		"visit_count":  1000,
	}

	// Необязательные поля добавляем только при выполнении условий
	if alias.Visit_count != 0 {
		fieldsMap["visit_count"] = alias.Visit_count
	}
	if alias.Expires_at != "" {
		fieldsMap["expires_at"] = alias.Expires_at
	}
	if alias.Title != "" {
		fieldsMap["title"] = alias.Title
	}
	if alias.Description != "" {
		fieldsMap["description"] = alias.Description
	}

	// Строим запрос из карты
	for columnName, value := range fieldsMap {
		addFieldToQuery(&columns, &placeholders, &values, columnName, value, &paramIndex)
	}

	query := fmt.Sprintf("INSERT INTO urls (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	_, err := storage.databasePool.Exec(context.Background(), query, values...)
	return err
}

func addFieldToQuery(columns *[]string, placeholders *[]string, values *[]interface{},
	columnName string, value interface{}, paramIndex *int) {
	*columns = append(*columns, columnName)
	*placeholders = append(*placeholders, fmt.Sprintf("$%d", *paramIndex))
	*values = append(*values, value)
	*paramIndex++
}
