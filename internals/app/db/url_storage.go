package db

import (
	"context"
	"encoding/json"
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
		return "", fmt.Errorf("Url for this alias not found err:%w", err)
	}
	return url, nil
}

func (storage *UrlStorage) GetUrlId(alias string) (int, error) {
	var url_id int
	query := "SELECT id FROM urls WHERE alias = $1"
	err := storage.databasePool.QueryRow(context.Background(), query, alias).Scan(&url_id)
	if err != nil {
		return -1, fmt.Errorf("Url for this alias not found err:%w", err)
	}
	return url_id, nil
}

// func (storage *UrlStorage) GetInfoAboutAlias(alias string) (*modals.UrlInfo, error) {
// 	var url modals.Url
// 	query := "SELECT * FROM urls WHERE alias = $1"
// 	err := storage.databasePool.QueryRow(context.Background(), query, alias).Scan(
// 		&url.Id,
// 		&url.Original_url,
// 		&url.Alias,
// 		&url.Created_at,
// 		&url.Expires_at,
// 		&url.Is_active,
// 		&url.Visit_count,
// 		&url.Title,
// 		&url.Description)
// 	if err != nil {
// 		return nil, fmt.Errorf("Can't get info about this alis err:%w", err)
// 	}
// 	return &url, nil
// }

func (storage *UrlStorage) GetInfoAboutAlias(alias string) (*modals.UrlInfo, error) {
	query := `
        SELECT 
            u.*,
            COALESCE(
                (SELECT json_agg(json_build_object(
                    'id', v.id,
                    'url_id', v.url_id,
                    'created_at', v.created_at
                )) 
                FROM visits v 
                WHERE v.url_id = u.id),
                '[]'::json
            ) as visits
        FROM urls u
        WHERE u.alias = $1
    `

	var url modals.UrlInfo
	var visitsJSON []byte

	err := storage.databasePool.QueryRow(context.Background(), query, alias).Scan(
		&url.Id,
		&url.Original_url,
		&url.Alias,
		&url.Created_at,
		&url.Expires_at,
		&url.Is_active,
		&url.Visit_count,
		&url.Title,
		&url.Description,
		&visitsJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get info about this alias err: %w", err)
	}
	var visits []modals.Visit
	if err := json.Unmarshal(visitsJSON, &visits); err != nil {
		return nil, fmt.Errorf("Can't unmarshal visits json: %w", err)
	}
	url.Visits = visits
	return &url, nil
}

func (storage *UrlStorage) CheckVisitCount(alias string) error {
	var count int
	query := "SELECT visit_count FROM urls WHERE alias = $1"
	err := storage.databasePool.QueryRow(context.Background(), query, alias).Scan(&count)
	if err != nil {
		return fmt.Errorf("Can't get visit count for this alias err:%w", err)
	}
	if count == 0 {
		return fmt.Errorf("This alias is unavailable, the number of visits has ended")
	}

	return nil
}

func (storage *UrlStorage) UpdateVisitCount(alias string) error {
	query := "UPDATE urls SET visit_count = visit_count - 1 WHERE alias = $1"
	_, err := storage.databasePool.Exec(context.Background(), query, alias)
	if err != nil {
		return fmt.Errorf("Can't update visit count for this alias err: %w", err)
	}
	var count int
	query = "SELECT visit_count FROM urls WHERE alias = $1"
	err = storage.databasePool.QueryRow(context.Background(), query, alias).Scan(&count)
	if count == 0 {
		query := "UPDATE urls SET is_active = FALSE WHERE alias = $1"
		_, err := storage.databasePool.Exec(context.Background(), query, alias)
		if err != nil {
			return fmt.Errorf("Can't update is_active for this alias err: %w", err)
		}
	}
	return nil
}

func (storage *UrlStorage) CheckExpiresAt(alias string) error {
	var expires_at *time.Time
	query := "SELECT expires_at FROM urls WHERE alias = $1"
	err := storage.databasePool.QueryRow(context.Background(), query, alias).Scan(&expires_at)
	if err != nil {
		return fmt.Errorf("Can't get expires_at time for this alias err:%w", err)
	}
	if expires_at != nil {
		if expires_at.Before(time.Now()) {
			return fmt.Errorf("This alias is unavailable, expired")
		}
	}

	return nil
}

func (storage *UrlStorage) DeleteAlias(alias string) error {
	query := "DELETE FROM urls WHERE alias = $1"
	_, err := storage.databasePool.Exec(context.Background(), query, alias)
	if err != nil {
		return fmt.Errorf("Can't delete alias err:%w", err)
	}
	return nil
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
	if alias.Expires_at != nil {
		fieldsMap["expires_at"] = &alias.Expires_at
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

func (storage *UrlStorage) GetAliases() ([]modals.Url, error) {
	var urls []modals.Url
	query := "SELECT * FROM urls"
	rows, err := storage.databasePool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Can't get all aliases err: %w", err)
	}

	for rows.Next() {
		var url modals.Url
		err := rows.Scan(
			&url.Id,
			&url.Original_url,
			&url.Alias,
			&url.Created_at,
			&url.Expires_at,
			&url.Is_active,
			&url.Visit_count,
			&url.Title,
			&url.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("Can't scan alias err: %w", err)
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (storage *UrlStorage) AddVisit(url_id int) error {
	query := "INSERT INTO visits (url_id, created_at) VALUES ($1, $2)"
	_, err := storage.databasePool.Exec(context.Background(), query, url_id, time.Now())
	return err
}

func addFieldToQuery(columns *[]string, placeholders *[]string, values *[]interface{},
	columnName string, value interface{}, paramIndex *int) {
	*columns = append(*columns, columnName)
	*placeholders = append(*placeholders, fmt.Sprintf("$%d", *paramIndex))
	*values = append(*values, value)
	*paramIndex++
}
