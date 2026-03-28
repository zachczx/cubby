package market

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MarketPrice struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	FamilyID  uuid.UUID  `json:"-" db:"family_id"`
	LoggedBy  *uuid.UUID `json:"loggedBy" db:"logged_by"`
	ItemName  string     `json:"itemName" db:"item_name"`
	Category  *string    `json:"category" db:"category"`
	Country   *string    `json:"country" db:"country"`
	Store     *string    `json:"store" db:"store"`
	Unit      *string    `json:"unit" db:"unit"`
	Quantity  *float64   `json:"quantity" db:"quantity"`
	Price     float64    `json:"price" db:"price"`
	IsPromo   bool       `json:"isPromo" db:"is_promo"`
	Remarks   *string    `json:"remarks" db:"remarks"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}

type MarketInsight struct {
	ItemName     string    `json:"itemName" db:"item_name"`
	Category     *string   `json:"category" db:"category"`
	Country      *string   `json:"country" db:"country"`
	LowestPrice  float64   `json:"lowestPrice" db:"lowest_price"`
	LowestUnit   *float64  `json:"lowestUnit" db:"lowest_unit"`
	LowestStore  *string   `json:"lowestStore" db:"lowest_store"`
	LowestDate   time.Time `json:"lowestDate" db:"lowest_date"`
	LatestPrice  float64   `json:"latestPrice" db:"latest_price"`
	LatestUnit   *float64  `json:"latestUnit" db:"latest_unit"`
	LatestStore  *string   `json:"latestStore" db:"latest_store"`
	LatestDate   time.Time `json:"latestDate" db:"latest_date"`
	DeltaPercent float64   `json:"deltaPercent" db:"delta_percent"`
}

type Input struct {
	ItemName string   `json:"itemName"`
	Category *string  `json:"category"`
	Country  *string  `json:"country"`
	Store    *string  `json:"store"`
	Unit     *string  `json:"unit"`
	Quantity *float64 `json:"quantity"`
	Price    float64  `json:"price"`
	IsPromo  bool     `json:"isPromo"`
	Remarks  *string  `json:"remarks"`
}

type UpsertResult struct {
	ID       uuid.UUID `json:"id"`
	IsUpdate bool      `json:"isUpdate"`
}

func LogPrice(db *sqlx.DB, p MarketPrice) (UpsertResult, error) {
	var result UpsertResult

	tx, err := db.Beginx()
	if err != nil {
		return result, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var existingID uuid.UUID
	checkQ := `SELECT id FROM market_prices
		WHERE family_id = $1
		AND LOWER(item_name) = LOWER($2)
		AND COALESCE(LOWER(store), '') = COALESCE(LOWER($3), '')
		AND price = $4
		AND DATE(created_at AT TIME ZONE 'UTC') = CURRENT_DATE
		LIMIT 1`

	err = tx.Get(&existingID, checkQ, p.FamilyID, p.ItemName, p.Store, p.Price)

	if err == nil {
		updateQ := `UPDATE market_prices SET
				quantity = $1, unit = $2, is_promo = $3, remarks = $4,
				logged_by = $5, updated_at = NOW()
			WHERE id = $6`
		if _, err := tx.Exec(updateQ, p.Quantity, p.Unit, p.IsPromo, p.Remarks, p.LoggedBy, existingID); err != nil {
			return result, fmt.Errorf("update duplicate: %w", err)
		}
		result.ID = existingID
		result.IsUpdate = true
	} else {
		insertQ := `INSERT INTO market_prices (
				family_id, logged_by, item_name, category, country, store, unit, quantity, price, is_promo, remarks, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
			RETURNING id`
		if err := tx.Get(&result.ID, insertQ,
			p.FamilyID, p.LoggedBy, p.ItemName, p.Category, p.Country,
			p.Store, p.Unit, p.Quantity, p.Price, p.IsPromo, p.Remarks,
		); err != nil {
			return result, fmt.Errorf("insert market price: %w", err)
		}
		result.IsUpdate = false
	}

	if err := tx.Commit(); err != nil {
		return result, fmt.Errorf("commit: %w", err)
	}

	return result, nil
}

func GetPrices(db *sqlx.DB, userID uuid.UUID) ([]MarketPrice, error) {
	var p []MarketPrice
	q := `SELECT mp.* FROM market_prices mp
			WHERE mp.family_id IN (
				SELECT family_id FROM families_users WHERE user_id = $1
				UNION
				SELECT id FROM families WHERE owner_id = $1
			)
			ORDER BY mp.created_at DESC`

	if err := db.Select(&p, q, userID); err != nil {
		return nil, fmt.Errorf("select market prices: %w", err)
	}

	if p == nil {
		p = []MarketPrice{}
	}

	return p, nil
}

type itemKey struct {
	name    string
	country string
}

func makeKey(name string, country *string) itemKey {
	k := itemKey{name: strings.ToLower(name)}
	if country != nil {
		k.country = strings.ToLower(*country)
	}
	return k
}

type latestRow struct {
	ItemName  string    `db:"item_name"`
	Category  *string   `db:"category"`
	Country   *string   `db:"country"`
	Price     float64   `db:"price"`
	UnitPrice float64   `db:"unit_price"`
	Store     *string   `db:"store"`
	CreatedAt time.Time `db:"created_at"`
}

type lowestRow struct {
	ItemName  string    `db:"item_name"`
	Country   *string   `db:"country"`
	Price     float64   `db:"price"`
	UnitPrice float64   `db:"unit_price"`
	Store     *string   `db:"store"`
	CreatedAt time.Time `db:"created_at"`
}

func getLatestPrices(db *sqlx.DB, userID uuid.UUID) ([]latestRow, error) {
	var rows []latestRow
	q := `SELECT DISTINCT ON (LOWER(item_name), LOWER(COALESCE(country, '')))
		item_name, category, country, price,
		COALESCE(price / NULLIF(quantity, 0), price) AS unit_price,
		store, created_at
	FROM market_prices
	WHERE family_id IN (
		SELECT family_id FROM families_users WHERE user_id = $1
		UNION
		SELECT id FROM families WHERE owner_id = $1
	)
	ORDER BY LOWER(item_name), LOWER(COALESCE(country, '')), created_at DESC`

	if err := db.Select(&rows, q, userID); err != nil {
		return nil, fmt.Errorf("select latest prices: %w", err)
	}
	return rows, nil
}

func getLowestPrices(db *sqlx.DB, userID uuid.UUID) ([]lowestRow, error) {
	var rows []lowestRow
	q := `SELECT DISTINCT ON (LOWER(item_name), LOWER(COALESCE(country, '')))
		item_name, country, price,
		COALESCE(price / NULLIF(quantity, 0), price) AS unit_price,
		store, created_at
	FROM market_prices
	WHERE family_id IN (
		SELECT family_id FROM families_users WHERE user_id = $1
		UNION
		SELECT id FROM families WHERE owner_id = $1
	)
	ORDER BY LOWER(item_name), LOWER(COALESCE(country, '')), unit_price ASC, created_at DESC`

	if err := db.Select(&rows, q, userID); err != nil {
		return nil, fmt.Errorf("select lowest prices: %w", err)
	}
	return rows, nil
}

func GetInsights(db *sqlx.DB, userID uuid.UUID) ([]MarketInsight, error) {
	latest, err := getLatestPrices(db, userID)
	if err != nil {
		return nil, err
	}

	lowest, err := getLowestPrices(db, userID)
	if err != nil {
		return nil, err
	}

	lowestMap := make(map[itemKey]lowestRow, len(lowest))
	for _, r := range lowest {
		lowestMap[makeKey(r.ItemName, r.Country)] = r
	}

	insights := make([]MarketInsight, 0, len(latest))
	for _, l := range latest {
		low, ok := lowestMap[makeKey(l.ItemName, l.Country)]
		if !ok {
			continue
		}

		var delta float64
		if low.UnitPrice != 0 {
			delta = (l.UnitPrice - low.UnitPrice) / low.UnitPrice * 100
			delta = float64(int(delta*10)) / 10 // round to 1 decimal
		}

		insights = append(insights, MarketInsight{
			ItemName:     l.ItemName,
			Category:     l.Category,
			Country:      l.Country,
			LowestPrice:  low.Price,
			LowestUnit:   &low.UnitPrice,
			LowestStore:  low.Store,
			LowestDate:   low.CreatedAt,
			LatestPrice:  l.Price,
			LatestUnit:   &l.UnitPrice,
			LatestStore:  l.Store,
			LatestDate:   l.CreatedAt,
			DeltaPercent: delta,
		})
	}

	if insights == nil {
		insights = []MarketInsight{}
	}

	return insights, nil
}

func DeletePrice(db *sqlx.DB, userID uuid.UUID, priceID uuid.UUID) error {
	q := `DELETE FROM market_prices
			WHERE id = $1
			AND family_id IN (
				SELECT family_id FROM families_users WHERE user_id = $2
				UNION
				SELECT id FROM families WHERE owner_id = $2
			)`

	if _, err := db.Exec(q, priceID, userID); err != nil {
		return fmt.Errorf("delete market price: %w", err)
	}

	return nil
}

func UpdatePrice(db *sqlx.DB, p MarketPrice, userID uuid.UUID) error {
	q := `UPDATE market_prices SET
			item_name = $1,
			category = $2,
			country = $3,
			store = $4,
			unit = $5,
			quantity = $6,
			price = $7,
			is_promo = $8,
			remarks = $9,
			updated_at = NOW()
		WHERE id = $10
		AND family_id IN (
			SELECT family_id FROM families_users WHERE user_id = $11
			UNION
			SELECT id FROM families WHERE owner_id = $11
		)`

	_, err := db.Exec(q,
		p.ItemName, p.Category, p.Country, p.Store, p.Unit,
		p.Quantity, p.Price, p.IsPromo, p.Remarks,
		p.ID, userID,
	)
	if err != nil {
		return fmt.Errorf("update market price: %w", err)
	}

	return nil
}
